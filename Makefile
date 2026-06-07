BINARY := nous
VERSION := 0.1.0

.PHONY: build test lint clean

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/nous

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY)
