FROM golang

WORKDIR /go/src/app

COPY . .

RUN go mod download

CMD ["make", "test-cmd"]

