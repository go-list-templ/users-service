FROM golang:1.24.4 AS build

WORKDIR /go/src/app

COPY . .

RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/app ./cmd

# todo add max min size image
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /go/bin/app /

CMD ["/app"]

