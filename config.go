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

// LoadConfig resolves the config path: -c path wins, then a positional arg,
// then ./gateproxy.json, then built-in defaults.
func LoadConfig() (*Config, error) {
path := "gateproxy.json"
args := os.Args[1:]
for i, a := range args {
if a == "-c" && i+1 < len(args) {
path = args[i+1]
args = append(args[:i], args[i+2:]...)
break
}
}
if len(args) > 0 {
path = args[0]
}

cfg := &Config{Listen: ":8080"}

raw, err := os.ReadFile(path)
switch {
case os.IsNotExist(err) && path == "gateproxy.json":
log.Printf("no config file at %s, using defaults", path)
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
cfg.Routes[i].RatePerMinute = 100
}
}
return cfg, nil
}

func defaultRoutes() []Route {
return []Route{
{Prefix: "/api/v1/", Upstream: "127.0.0.1:9001", RatePerMinute: 100},
}
}
