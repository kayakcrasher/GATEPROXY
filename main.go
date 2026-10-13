package main

import (
"log"
"net"
)

func main() {
cfg, err := LoadConfig()
if err != nil {
log.Fatalf("config error: %v", err)
}
routes = cfg.Routes

listener, err := net.Listen("tcp", cfg.Listen)
if err != nil {
log.Fatalf("failed to listen on %s: %v", cfg.Listen, err)
}
defer listener.Close()

log.Printf("GATEPROXY listening on %s with %d route(s)", cfg.Listen, len(routes))
for _, r := range routes {
log.Printf("  %s -> %s (%d req/min per IP)", r.Prefix, r.Upstream, r.RatePerMinute)
}

for {
conn, err := listener.Accept()
if err != nil {
log.Printf("accept error: %v", err)
continue
}
go handleClient(conn)
}
}
