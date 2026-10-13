# ---- Build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /out/gateproxy .

# ---- Runtime stage ----
FROM alpine:3.19
COPY --from=build /out/gateproxy /usr/local/bin/gateproxy
COPY gateproxy.json /etc/gateproxy/gateproxy.json
EXPOSE 8080
ENTRYPOINT ["gateproxy"]
CMD ["-c", "/etc/gateproxy/gateproxy.json"]
