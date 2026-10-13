BINARY  := gateproxy
VERSION := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

build:
go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

run: build
./$(BINARY)

vet:
go vet ./...

test:
go test ./...

clean:
rm -f $(BINARY)

docker-build:
docker build -t gateproxy:$(VERSION) .

.PHONY: build run vet test clean docker-build
