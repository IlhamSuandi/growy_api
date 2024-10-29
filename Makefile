run: build
	@./bin/api

air:
	@air

run-dev:
	@APP_ENV=development go run cmd/server/main.go

run-prod:
	@APP_ENV=production go run cmd/server/main.go

run-staging:
	@APP_ENV=staging go run cmd/server/main.go

build:
	@go build -C cmd/server/ -o ../../bin/api 

seeder:
	@go run ./cmd/seeder

swagger-init:
	@swag init -d ./ -g ./cmd/server/main.go

tests:
	@go test -v ./test/...

tests-%:
	@go test -v ./test/... -run=$(shell echo $* | sed 's/_/./g')

migration:
	@migrate create -ext sql -dir database/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down

