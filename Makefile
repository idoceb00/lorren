.PHONY: build install test lint fmt

build:
	go build -o bin/lorren ./cmd/lorren

install:
	go install ./cmd/lorren

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofmt -l .