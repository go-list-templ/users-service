FROM golang:1.24.4

COPY . .

RUN go mod download

CMD ["go", "test", "-v", "-count=1", "./internal/..."]

