FROM golang:1.26-alpine3.23 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/app /
COPY --from=build /app/migrations /migrations

ENTRYPOINT ["/app"]

