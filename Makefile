.PHONY: cli
cli:
		CGO_ENABLED=0 go build -o . cmd/cli/main.go

run-cli:
		go run cmd/cli/main.go init