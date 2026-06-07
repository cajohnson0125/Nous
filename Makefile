.PHONY: build test lint vet clean

BINARY := nous
BUILD_FLAGS := -ldflags="-s -w"

build:
	go build $(BUILD_FLAGS) ./cmd/$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

lint: vet
	golangci-lint run

clean:
	rm -f $(BINARY)
