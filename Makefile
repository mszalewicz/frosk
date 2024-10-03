.PHONY: run run_linux build build_linux

run:
	@go run cmd/main.go

run_linux:
	@go run --tags nowayland cmd/main.go

build:
	@go build -ldflags="-s -w" -o bin/frosk cmd/main.go

build_linux:
	@go build --tags nowayland -ldflags="-s -w" -o bin/frosk cmd/main.go
