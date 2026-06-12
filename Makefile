VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/edgedelta/edx/internal/cli.Version=$(VERSION)

.PHONY: build install test vet lint clean

build:
	go build -ldflags '$(LDFLAGS)' -o bin/edx .

install:
	go install -ldflags '$(LDFLAGS)' .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf bin
