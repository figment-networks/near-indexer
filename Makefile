.PHONY: build test fmt

build:
	go build

test:
	go test -cover -race ./...

fmt:
	go fmt ./...