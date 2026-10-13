# GATEPROXY [![CI](https://github.com/kayakcrasher/GATEPROXY/actions/workflows/ci.yml/badge.svg)](https://github.com/kayakcrasher/GATEPROXY/actions/workflows/ci.yml)

A high-performance, lightweight API Gateway and reverse proxy built from scratch in Go — raw TCP sockets, zero framework bloat, zero external runtime dependencies.

## Features
- **Raw TCP Sockets & Upstream Forwarding:** Concurrent request handling with streaming body forwarding and hop-by-hop header stripping.
- **Path-Prefix Routing:** Configurable via JSON route tables.
- **Token-Bucket Rate Limiting:** Per-IP sliding-window rate defense returning `429 Too Many Requests`.
- **Active Health Checks:** Background monitoring that marks dead upstreams down and fails fast with `503 Service Unavailable`.
- **Structured JSON Logging:** One parseable JSON object per request tracking latency, status, bytes, and client IPs.
- **Container Ready:** Multi-stage Dockerfile and a clean `Makefile`.

## Getting Started

### 1. Build
```bash
make build
