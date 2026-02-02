FROM golang:1.25

WORKDIR /go/src/app

COPY . .

RUN go mod download

CMD ["make", "test-cmd"]

