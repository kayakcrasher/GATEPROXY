package main

import (
"encoding/json"
"os"
"sync"
"time"
)

// LogEntry is one JSON line of structured output.
type LogEntry struct {
Time     string `json:"time"`
Level    string `json:"level"`
Msg      string `json:"msg"`
Method   string `json:"method,omitempty"`
Path     string `json:"path,omitempty"`
Status   int    `json:"status,omitempty"`
ClientIP string `json:"client_ip,omitempty"`
Upstream string `json:"upstream,omitempty"`
Latency  string `json:"latency,omitempty"`
Bytes    int64  `json:"bytes,omitempty"`
}

var logMu sync.Mutex

func logJSON(e LogEntry) {
e.Time = time.Now().UTC().Format(time.RFC3339)
raw, err := json.Marshal(e)
if err != nil {
return
}
logMu.Lock()
os.Stdout.Write(append(raw, '\n'))
logMu.Unlock()
}

func logEvent(level, msg string) {
logJSON(LogEntry{Level: level, Msg: msg})
}

func logRequest(clientIP, method, path string, status int, bytes int64, latency time.Duration, upstream string) {
logJSON(LogEntry{
Level:    "info",
Msg:      "request",
Method:   method,
Path:     path,
Status:   status,
ClientIP: clientIP,
Upstream: upstream,
Latency:  latency.String(),
Bytes:    bytes,
})
}
