test:
	go test -race -short ./...

integration-tests:
	go run github.com/onsi/ginkgo/v2/ginkgo@latest -p -race ./integration_tests/...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v --timeout 5m
