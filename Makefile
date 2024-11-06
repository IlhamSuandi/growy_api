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
	@read -p "Enter migration name: " name; \
	cd database/migrations; goose create $$name sql

migration-status:
	@goose -dir database/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" status

migrate-up:
	@goose -dir database/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	@goose -dir database/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down
