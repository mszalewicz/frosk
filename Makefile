.PHONY: run

run:
	rm -f sql/application.sqlite
	rm -f cmd/log
	@go run --tags nowayland cmd/main.go
