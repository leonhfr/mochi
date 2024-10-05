.PHONY: default
default: lint test

.PHONY: build
build:
	goreleaser release --snapshot --clean

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -v -coverprofile=coverage.out ./...

.PHONY: coverage-html
coverage-html: test
	go tool cover -html=coverage.out
