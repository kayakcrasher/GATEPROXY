package main

import (
"encoding/json"
"fmt"
"log"
"os"
)

// Config is the top-level structure of gateproxy.json.
type Config struct {
Listen string  `json:"listen"`
Routes []Route `json:"routes"`
}

// LoadConfig reads gateproxy.json (or the file passed as the first CLI arg),
// validates it, and falls back to sensible defaults when fields are missing.
func LoadConfig() (*Config, error) {
path := "gateproxy.json"
if len(os.Args) > 1 {
path = os.Args[1]
}

cfg := &Config{Listen: ":8080"}

raw, err := os.ReadFile(path)
switch {
case os.IsNotExist(err) && len(os.Args) <= 1:
// No config file and no explicit path: use defaults, log a warning.
log.Printf("no config file found at %s, using defaults", path)
cfg.Routes = defaultRoutes()
return cfg, nil
case err != nil:
return nil, fmt.Errorf("reading config: %w", err)
}

if err := json.Unmarshal(raw, cfg); err != nil {
return nil, fmt.Errorf("parsing %s: %w", path, err)
}

if cfg.Listen == "" {
cfg.Listen = ":8080"
}
if len(cfg.Routes) == 0 {
return nil, fmt.Errorf("%s: no routes defined", path)
}
for i, r := range cfg.Routes {
if r.Prefix == "" || r.Upstream == "" {
return nil, fmt.Errorf("%s: route %d missing prefix or upstream", path, i)
}
if r.RatePerMinute <= 0 {
cfg.Routes[i].RatePerMinute = 100 // safe default
}
}
return cfg, nil
}

func defaultRoutes() []Route {
return []Route{
{Prefix: "/api/v1/", Upstream: "127.0.0.1:9001", RatePerMinute: 100},
}
}
