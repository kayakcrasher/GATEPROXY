package main

import (
"sync"
"time"
)

type bucket struct {
tokens   float64
lastFill time.Time
}

// RateLimiter is a token-bucket limiter keyed by client IP.
type RateLimiter struct {
mu      sync.Mutex
buckets map[string]*bucket
rate    float64 // tokens refilled per second
burst   float64 // max bucket size
}

func NewRateLimiter(perMinute int) *RateLimiter {
return &RateLimiter{
buckets: make(map[string]*bucket),
rate:    float64(perMinute) / 60.0,
burst:   float64(perMinute),
}
}

func (rl *RateLimiter) Allow(key string) bool {
rl.mu.Lock()
defer rl.mu.Unlock()

now := time.Now()
b, ok := rl.buckets[key]
if !ok {
b = &bucket{tokens: rl.burst, lastFill: now}
rl.buckets[key] = b
}

// Refill based on elapsed time
b.tokens += now.Sub(b.lastFill).Seconds() * rl.rate
if b.tokens > rl.burst {
b.tokens = rl.burst
}
b.lastFill = now

if b.tokens >= 1 {
b.tokens--
return true
}
return false
}

// limiters keeps one RateLimiter per configured rate so routes
// with the same limit share nothing but code.
var limiters = struct {
sync.Mutex
m map[int]*RateLimiter
}{m: make(map[int]*RateLimiter)}

func limiterFor(perMinute int) *RateLimiter {
limiters.Lock()
defer limiters.Unlock()
if l, ok := limiters.m[perMinute]; ok {
return l
}
l := NewRateLimiter(perMinute)
limiters.m[perMinute] = l
return l
}
