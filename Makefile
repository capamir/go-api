build:
	@go build -o bin/app cmd/main.go

run: build
	@./bin/app

test:
	@go test -v ./...

dev:
	@go run cmd/main.go