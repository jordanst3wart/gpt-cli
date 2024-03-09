.PHONY: build test clean

build:
	go build -o dist/ -v ./...

test:
	go test -v ./...

clean:
	rm -f $(BIN)

