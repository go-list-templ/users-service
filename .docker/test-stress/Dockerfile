FROM alpine/git:latest AS loader
WORKDIR /app

RUN git clone https://github.com/go-list-templ/proto.git

FROM grafana/k6:latest

WORKDIR /src

COPY ./test/stress /src

COPY --from=loader /app/proto /src/proto

ENTRYPOINT ["k6", "run", "/src/main.js", "--compatibility-mode=base"]