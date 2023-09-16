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
	go test -race -v ./...
