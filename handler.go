package main

import (
"bufio"
"fmt"
"io"
"net"
"strconv"
"strings"
"time"
)

type countingWriter struct {
conn  net.Conn
count int64
}

func (cw *countingWriter) Write(p []byte) (int, error) {
n, err := cw.conn.Write(p)
cw.count += int64(n)
return n, err
}

func handleClient(client net.Conn) {
defer client.Close()
start := time.Now()
_ = client.SetDeadline(start.Add(30 * time.Second))
clientIP := ipOnly(client.RemoteAddr().String())

reader := bufio.NewReader(client)

// --- Request line ---
requestLine, err := reader.ReadString('\n')
if err != nil {
return
}
parts := strings.Fields(strings.TrimSpace(requestLine))
method, path := "-", "-"
if len(parts) >= 3 {
method, path = parts[0], parts[1]
} else {
writeSimpleResponse(client, 400, "Bad Request")
logRequest(clientIP, method, path, 400, 0, time.Since(start), "")
return
}

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
logRequest(clientIP, method, path, 404, 0, time.Since(start), "")
return
}

// --- Rate limiting ---
if !limiterFor(route.RatePerMinute).Allow(clientIP) {
writeSimpleResponse(client, 429, "Too Many Requests")
logRequest(clientIP, method, path, 429, 0, time.Since(start), route.Upstream)
return
}

// --- Health check: known-dead upstream -> 503 ---
if IsDown(route.Upstream) {
writeSimpleResponse(client, 503, "Service Unavailable")
logRequest(clientIP, method, path, 503, 0, time.Since(start), route.Upstream)
return
}

// --- Dial upstream ---
backend, err := net.DialTimeout("tcp", route.Upstream, 5*time.Second)
if err != nil {
MarkDown(route.Upstream)
writeSimpleResponse(client, 502, "Bad Gateway")
logRequest(clientIP, method, path, 502, 0, time.Since(start), route.Upstream)
return
}
defer backend.Close()
_ = backend.SetDeadline(start.Add(30 * time.Second))

// --- Forward request ---
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
fmt.Fprintf(backend, "X-Forwarded-For: %s\r\nConnection: close\r\n\r\n", clientIP)

if contentLength > 0 {
if _, err := io.CopyN(backend, reader, contentLength); err != nil {
logEvent("warn", "body forward error: "+err.Error())
return
}
}

// --- Read backend status line, then stream the rest ---
breader := bufio.NewReader(backend)
statusLine, err := breader.ReadString('\n')
if err != nil {
MarkDown(route.Upstream)
writeSimpleResponse(client, 502, "Bad Gateway")
logRequest(clientIP, method, path, 502, 0, time.Since(start), route.Upstream)
return
}
status := parseStatus(statusLine)

cw := &countingWriter{conn: client}
cw.Write([]byte(statusLine))
io.Copy(cw, breader)

logRequest(clientIP, method, path, status, cw.count, time.Since(start), route.Upstream)
}

func parseStatus(line string) int {
fields := strings.Fields(strings.TrimSpace(line))
if len(fields) >= 2 {
if code, err := strconv.Atoi(fields[1]); err == nil {
return code
}
}
return 502
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
