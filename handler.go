package main

import (
"bufio"
"fmt"
"io"
"log"
"net"
"strconv"
"strings"
"time"
)

func handleClient(client net.Conn) {
defer client.Close()
_ = client.SetDeadline(time.Now().Add(30 * time.Second))

reader := bufio.NewReader(client)

// --- Request line: METHOD /path HTTP/1.1 ---
requestLine, err := reader.ReadString('\n')
if err != nil {
return
}
parts := strings.Fields(strings.TrimSpace(requestLine))
if len(parts) < 3 {
writeSimpleResponse(client, 400, "Bad Request")
return
}
method, path := parts[0], parts[1]

// --- Headers ---
var contentLength int64
var headerLines []string
for {
line, err := reader.ReadString('\n')
if err != nil {
return
}
if line == "\r\n" || line == "\n" {
break
}
clean := strings.TrimRight(line, "\r\n")
headerLines = append(headerLines, clean)
if idx := strings.Index(clean, ":"); idx >= 0 {
if strings.EqualFold(strings.TrimSpace(clean[:idx]), "Content-Length") {
contentLength, _ = strconv.ParseInt(strings.TrimSpace(clean[idx+1:]), 10, 64)
}
}
}

// --- Routing ---
route, ok := MatchRoute(path)
if !ok {
writeSimpleResponse(client, 404, "No Route")
return
}

// --- Rate limiting (per client IP) ---
clientIP := ipOnly(client.RemoteAddr().String())
if !limiterFor(route.RatePerMinute).Allow(clientIP) {
log.Printf("RATE LIMITED %s %s from %s", method, path, clientIP)
writeSimpleResponse(client, 429, "Too Many Requests")
return
}

// --- Connect upstream ---
backend, err := net.DialTimeout("tcp", route.Upstream, 5*time.Second)
if err != nil {
log.Printf("upstream unreachable: %v", err)
writeSimpleResponse(client, 502, "Bad Gateway")
return
}
defer backend.Close()
_ = backend.SetDeadline(time.Now().Add(30 * time.Second))

// --- Forward request line + headers (strip hop-by-hop) ---
fmt.Fprintf(backend, "%s %s HTTP/1.1\r\n", method, path)
for _, h := range headerLines {
idx := strings.Index(h, ":")
if idx < 0 {
continue
}
lower := strings.ToLower(strings.TrimSpace(h[:idx]))
if lower == "connection" || lower == "keep-alive" || strings.HasPrefix(lower, "proxy-") {
continue
}
fmt.Fprintf(backend, "%s\r\n", h)
}
fmt.Fprintf(backend, "X-Forwarded-For: %s\r\n", clientIP)
fmt.Fprintf(backend, "Connection: close\r\n\r\n")

// --- Forward the body (Content-Length bytes), streamed ---
if contentLength > 0 {
if _, err := io.CopyN(backend, reader, contentLength); err != nil {
log.Printf("body forward error: %v", err)
return
}
}

log.Printf("%s %s -> %s (body=%d bytes)", method, path, route.Upstream, contentLength)

// --- Stream response back to client ---
io.Copy(client, backend)
}

func writeSimpleResponse(conn net.Conn, status int, reason string) {
body := reason + "\n"
fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n%s",
status, reason, len(body), body)
}

func ipOnly(addr string) string {
if host, _, err := net.SplitHostPort(addr); err == nil {
return host
}
return addr
}
