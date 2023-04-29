test:
	go test -race -short ./...

all-tests:
	go test -race ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v --timeout 5m
