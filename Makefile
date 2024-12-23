.PHONY: default
default: tidy lint test

.PHONY: build
build:
	goreleaser release --snapshot --clean

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -coverprofile=coverage.out ./...

.PHONY: coverage-html
coverage-html: test
	go tool cover -html=coverage.out

.PHONY: tidy
tidy:
	go mod tidy
