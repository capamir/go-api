build:
	@go build -o bin/app cmd/main.go

run: build
	@./bin/app

test:
	@go test -v ./...

dev:
	@go run cmd/main.go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
