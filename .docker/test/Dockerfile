FROM golang:1.26

WORKDIR /go/src/app

COPY . .

RUN go mod download

CMD ["make", "test-cmd"]

