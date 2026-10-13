package main

import "strings"

// Route maps a path prefix to an upstream backend and a per-IP rate limit.
type Route struct {
Prefix        string `json:"prefix"`
Upstream      string `json:"upstream"`
RatePerMinute int    `json:"rate_per_minute"`
}

// routes is populated by LoadConfig at startup.
var routes []Route

// MatchRoute returns the first route whose prefix matches the path.
func MatchRoute(path string) (Route, bool) {
for _, r := range routes {
if strings.HasPrefix(path, r.Prefix) {
return r, true
}
}
return Route{}, false
}
