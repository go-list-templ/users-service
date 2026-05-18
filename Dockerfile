FROM alpine:3.23 AS tools
RUN apk add --no-cache curl tar && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.1/migrate.linux-amd64.tar.gz | tar xvz -C /tmp

FROM golang:1.26-alpine3.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/app /
COPY --from=build /app/migrations /migrations/schemes
COPY --from=tools /tmp/migrate /migrations/migrate

EXPOSE 8080 8081

ENTRYPOINT ["/app"]