package main

import (
"testing"
)

func TestRateLimiter_Allow(t *testing.T) {
// 60 tokens per minute = 1 token per second. Burst = 60.
rl := NewRateLimiter(60)

ip := "192.168.1.50"

// Should allow the first burst of requests
for i := 0; i < 60; i++ {
if !rl.Allow(ip) {
t.Fatalf("expected request %d to be allowed under burst limit", i+1)
}
}

// 61st request should be blocked immediately (burst exhausted)
if rl.Allow(ip) {
t.Fatalf("expected 61st request to be rate limited")
}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
rl := NewRateLimiter(10)

ip1 := "10.0.0.1"
ip2 := "10.0.0.2"

// Exhaust IP1
for i := 0; i < 10; i++ {
rl.Allow(ip1)
}

if rl.Allow(ip1) {
t.Fatalf("IP1 should be rate limited")
}

// IP2 should still have its full bucket
if !rl.Allow(ip2) {
t.Fatalf("IP2 should not be affected by IP1's traffic")
}
}
