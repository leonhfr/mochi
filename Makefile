.PHONY: default
default: lint test

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -v -coverprofile=coverage.out -coverpkg=github.com/leonhfr/mochi/... ./...

.PHONY: coverage-html
coverage-html: test
	go tool cover -html=coverage.out
