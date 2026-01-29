FROM golang:1.24.4

WORKDIR /go/src/app

COPY . .

RUN go mod download

CMD ["make", "test-cmd"]

