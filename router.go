package main

import "strings"

// Route maps a path prefix to an upstream backend and a per-IP rate limit.
type Route struct {
Prefix        string
Upstream      string
RatePerMinute int
}

var routes = []Route{
{Prefix: "/api/v1/", Upstream: "127.0.0.1:9001", RatePerMinute: 100},
{Prefix: "/api/v2/", Upstream: "127.0.0.1:9002", RatePerMinute: 10},
}

// MatchRoute returns the first route whose prefix matches the path.
func MatchRoute(path string) (Route, bool) {
for _, r := range routes {
if strings.HasPrefix(path, r.Prefix) {
return r, true
}
}
return Route{}, false
}
