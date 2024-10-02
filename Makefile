.PHONY: run build

run:
	@go run --tags nowayland cmd/main.go

build:
	@go build --tags nowayland -ldflags="-s -w" -o bin/frosk cmd/main.go
