code-lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.9.0 golangci-lint run --config=lint.yml

lint-fix:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v2.9.0 golangci-lint run --config=lint.yml --fix

test-cmd:
	go test -v -race -count=1 ./internal/...

test-coverage-cmd:
	go install github.com/vladopajic/go-test-coverage/v2@latest
	go test ./internal/... -coverprofile=./cover.out -covermode=atomic
	go-test-coverage --config=./.docker/test-coverage/conf.yml

test-integration-cmd:
	go test -v -count=1 ./test/intergration/...