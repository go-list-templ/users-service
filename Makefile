include .env

DB_URL := ${DB_DRIVER}://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_DATABASE}?sslmode=disable
MIGRATION_PATH := $$(pwd)/internal/adapter/persistence/postgres/migrations

start: build migrate

build:
	docker compose --env-file .env up -d --build --remove-orphans

migrate:
	docker run --rm --network=host -v "$(MIGRATION_PATH)":/db/migrations -e DATABASE_URL="$(DB_URL)" amacneil/dbmate:2.28 --migrations-table=migrations --wait up

migrate-rollback:
	docker run --rm --network=host -v "$(MIGRATION_PATH)":/db/migrations -e DATABASE_URL="$(DB_URL)" amacneil/dbmate:2.28 --migrations-table=migrations --wait rollback

log:
	docker logs -f --tail 100 app.${APP_NAME}

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run --config=.docker/lint/conf.yml

.PHONY: test
test:
	go generate ./...
	docker build -f ./.docker/test/Dockerfile -t go-test .
	docker run --rm go-test

test-cmd:
	go test -v -count=1 ./internal/...

test-coverage:
	go generate ./...
	docker build -f ./.docker/coverage/Dockerfile -t go-test-coverage .
	docker run --rm go-test-coverage

test-coverage-cmd:
	go install github.com/vladopajic/go-test-coverage/v2@latest
	go test ./internal/... -coverprofile=./cover.out -covermode=atomic
	go-test-coverage --config=./.docker/coverage/conf.yml