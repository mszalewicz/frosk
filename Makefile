.PHONY: run build

run:
	#rm -f sql/application.sqlite
	#rm -f cmd/log
	@go run --tags nowayland cmd/main.go

build:
	@go build --tags nowayland -ldflags="-s -w" -o bin/frosk cmd/main.go
