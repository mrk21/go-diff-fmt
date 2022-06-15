.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags '-extldflags=-static' -o bin/godiff cmd/godiff/main.go

.PHONY: test
test:
	go test -v ./difffmt
