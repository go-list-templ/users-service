include .env

DB_URL := ${DB_DRIVER}://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_DATABASE}?sslmode=disable
MIGRATION_PATH := $$(pwd)/internal/adapter/persistence/postgres/migrations

start: build migrate

lint: docker-lint go-lint

lint-fix: go-lint-fix

build:
	docker compose --env-file .env up -d --build --remove-orphans

migrate:
	docker run --rm --network=host -v "$(MIGRATION_PATH)":/db/migrations -e DATABASE_URL="$(DB_URL)" amacneil/dbmate:2.28 --migrations-table=migrations --wait up

migrate-rollback:
	docker run --rm --network=host -v "$(MIGRATION_PATH)":/db/migrations -e DATABASE_URL="$(DB_URL)" amacneil/dbmate:2.28 --migrations-table=migrations --wait rollback

log:
	docker logs -f --tail 100 app.${APP_NAME}

go-lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run --config=.docker/golangci/.golangci.yml

go-lint-fix:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run --fix --config=.docker/golangci/.golangci.yml

docker-lint:
	@for f in $(shell find . -name 'Dockerfile'); do \
    	  echo "Lint $$f"; \
    	  docker run --rm -i \
    	    -v "$$PWD":/src \
    	    -v ./.docker/hadolint/hadolint.yaml:/.config/hadolint.yaml \
    	    --workdir /src \
    	    hadolint/hadolint hadolint "$${f#./}"; \
    	done
