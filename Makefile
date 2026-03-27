.PHONY: run build test lint ci

run:
	go run ./cmd/main.go

build:
	CGO_ENABLED=0 go build -o bin/gateway ./cmd/main.go

test:
	go test -race ./...

lint:
	golangci-lint run ./...

ci: lint test build

docker-build:
	docker build -t portfolio-gateway .
