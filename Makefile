.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/godiff cmd/godiff/main.go

.PHONY: test
test:
	go test -v ./difffmt
