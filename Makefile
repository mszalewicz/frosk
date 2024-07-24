.PHONY: run

run:
	rm -f sql/application.sqlite
	rm -f cmd/log
	@go run cmd/main.go
