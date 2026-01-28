FROM golang:1.24.4

WORKDIR /go/src/app

COPY . .

RUN go mod download && \
    go install github.com/vladopajic/go-test-coverage/v2@v2.18.3

CMD ["make", "test-coverage-cmd"]

