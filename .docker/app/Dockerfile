FROM golang:1.26-alpine3.23 AS build

WORKDIR /go/src/app

COPY . .

RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /go/bin/app /

CMD ["/app"]

