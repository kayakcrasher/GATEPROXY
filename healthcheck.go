package main

import (
"net"
"sync"
"time"
)

type healthState struct {
mu   sync.Mutex
down map[string]bool
}

var health = &healthState{down: make(map[string]bool{})}

// StartHealthChecks pings every upstream on a ticker until the process exits.
func StartHealthChecks(upstreams []string, interval time.Duration) {
health.mu.Lock()
for _, u := range upstreams {
health.down[u] = false
}
health.mu.Unlock()

go func() {
ticker := time.NewTicker(interval)
defer ticker.Stop()
for range ticker.C {
for _, u := range upstreams {
go pingUpstream(u)
}
}
}()
}

func pingUpstream(u string) {
conn, err := net.DialTimeout("tcp", u, 2*time.Second)
wasDown := IsDown(u)
if err != nil {
health.mu.Lock()
health.down[u] = true
health.mu.Unlock()
if !wasDown {
logJSON(LogEntry{Level: "warn", Msg: "upstream down", Upstream: u})
}
return
}
conn.Close()
health.mu.Lock()
health.down[u] = false
health.mu.Unlock()
if wasDown {
logJSON(LogEntry{Level: "info", Msg: "upstream recovered", Upstream: u})
}
}

// IsDown reports whether an upstream is currently marked unhealthy.
func IsDown(u string) bool {
health.mu.Lock()
defer health.mu.Unlock()
return health.down[u]
}

// MarkDown flags an upstream unhealthy immediately (e.g. a live dial failed).
func MarkDown(u string) {
health.mu.Lock()
health.down[u] = true
health.mu.Unlock()
}
