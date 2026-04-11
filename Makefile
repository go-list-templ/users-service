lint: lint-code helm-lint

helm-lint:
	werf lint

lint-code:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.9.0 golangci-lint run --config=lint.yml

lint-fix:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.9.0 golangci-lint run --config=lint.yml --fix

unit-test:
	go test -v -race -count=1 ./internal/...

coverage-test:
	go install github.com/vladopajic/go-test-coverage/v2@latest
	go test ./internal/... -coverprofile=./cover.out -covermode=atomic
	go-test-coverage --config=coverage.yml

integration-test:
	go test -v -count=1 ./test/intergration/...