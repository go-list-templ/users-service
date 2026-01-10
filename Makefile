include .env

DB_URL=${DB_DRIVER}://${DB_USERNAME}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_DATABASE}?sslmode=disable

start: build migrate

lint: docker-lint go-lint

lint-fix: go-lint-fix

build:
	docker compose --env-file .env up -d --build --remove-orphans

docker-lint:
	docker run --rm -i -v ./hadolint.yaml:/.config/hadolint.yaml hadolint/hadolint < ./Dockerfile

go-lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run

go-lint-fix:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run --fix

migrate:
	docker run --rm --network=host -v $$(pwd)/db/migrations:/db/migrations -e DATABASE_URL="$(DB_URL)" amacneil/dbmate:2.28 up

migrate-rollback:
	docker run --rm --network=host -v $$(pwd)/db/migrations:/db/migrations -e DATABASE_URL="$(DB_URL)" amacneil/dbmate:2.28 rollback

log:
	docker logs -f --tail 10 app.${APP_NAME}
