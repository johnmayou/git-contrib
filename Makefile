.PHONY: all build test fmt lint tidy clean

all: build

build:
	@go build -o app.exe .

test:
	@go test ./... -v -count=1

fmt:
	@goimports -w .
	@gofmt -s -w .

lint:
	@golangci-lint run ./...

tidy:
	@go mod tidy
	@go mod verify

clean:
	@rm app.exe