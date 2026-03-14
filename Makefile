include .env

COMPOSE_CMD := docker compose -p ${APP_NAME} --env-file .env
COMPOSE_TEST_CMD := docker compose -p ${APP_NAME}_tests --env-file .env -f docker-compose.yml -f .docker/test-integration/docker-compose.yml
COMPOSE_TEST_STRESS_CMD := docker compose -p ${APP_NAME}_tests_stress --env-file .env -f docker-compose.yml -f .docker/test-stress/docker-compose.yml

init: build
lint: docker-lint code-lint

build:
	$(COMPOSE_TEST_CMD) down -v
	$(COMPOSE_TEST_STRESS_CMD) down -v
	$(COMPOSE_CMD) up -d --build --remove-orphans

down:
	$(COMPOSE_CMD) down -v

log:
	docker logs -f --tail 100 app.${APP_NAME}

code-lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.9.0 golangci-lint run --config=.docker/lint/conf.yml

lint-fix:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.9.0 golangci-lint run --config=.docker/lint/conf.yml --fix

docker-lint:
	docker run --rm -i hadolint/hadolint < .docker/app/Dockerfile

.PHONY: test
test:
	go generate ./...
	docker build -f ./.docker/test/Dockerfile -t go-test .
	docker run --rm go-test

test-cmd:
	go test -v -race -count=1 ./internal/...

test-coverage:
	go generate ./...
	docker build -f ./.docker/test-coverage/Dockerfile -t go-test-coverage .
	docker run --rm go-test-coverage

test-coverage-cmd:
	go install github.com/vladopajic/go-test-coverage/v2@latest
	go test ./internal/... -coverprofile=./cover.out -covermode=atomic
	go-test-coverage --config=./.docker/test-coverage/conf.yml

test-integration:
	$(COMPOSE_CMD) down
	$(COMPOSE_TEST_CMD) down -v
	$(COMPOSE_TEST_STRESS_CMD) down -v
	$(COMPOSE_TEST_CMD) up --build -V --abort-on-container-exit --exit-code-from test --attach test

test-integration-cmd:
	go test -v -count=1 ./test/intergration/...

test-stress:
	$(COMPOSE_CMD) down
	$(COMPOSE_TEST_CMD) down -v
	$(COMPOSE_TEST_STRESS_CMD) down -v
	$(COMPOSE_TEST_STRESS_CMD) up --build -V --abort-on-container-exit --exit-code-from test-stress --attach test-stress