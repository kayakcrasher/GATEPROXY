package main

import (
"log"
"net"
"os"
"time"
)

// version is injected at build time via -ldflags.
var version = "dev"

func main() {
if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
logEvent("info", "gateproxy "+version)
return
}

cfg, err := LoadConfig()
if err != nil {
log.Fatalf("config error: %v", err)
}
routes = cfg.Routes

seen := make(map[string]bool)
var upstreams []string
for _, r := range routes {
if !seen[r.Upstream] {
seen[r.Upstream] = true
upstreams = append(upstreams, r.Upstream)
}
}
StartHealthChecks(upstreams, 10*time.Second)

listener, err := net.Listen("tcp", cfg.Listen)
if err != nil {
log.Fatalf("failed to listen on %s: %v", cfg.Listen, err)
}
defer listener.Close()

logJSON(LogEntry{Level: "info", Msg: "gateproxy listening", Upstream: cfg.Listen, Path: "version=" + version})
for _, r := range routes {
logJSON(LogEntry{Level: "info", Msg: "route", Path: r.Prefix, Upstream: r.Upstream, Status: r.RatePerMinute})
}

for {
conn, err := listener.Accept()
if err != nil {
logEvent("warn", "accept error: "+err.Error())
continue
}
go handleClient(conn)
}
}
