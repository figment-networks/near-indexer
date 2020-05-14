.PHONY: build test

build:
	go build

test:
	go test -cover -race ./...

fmt:
	go fmt ./...