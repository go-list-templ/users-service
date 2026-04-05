FROM golang:1.26-alpine3.23

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -v -race -count=1 ./internal/...

