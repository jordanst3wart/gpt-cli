.PHONY: build test clean

# The binary to build (just the basename).
BIN := bin

# This version-strategy uses git tags to set the version string
VERSION ?= $(shell git describe --tags --always --dirty)

build:
	go build -o $(BIN) -v ./...

test:
	go test -v ./...

clean:
	rm -f $(BIN)

