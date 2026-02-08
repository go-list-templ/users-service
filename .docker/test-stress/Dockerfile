FROM grafana/k6:latest

COPY ./test/stress /src
WORKDIR /src

ENTRYPOINT ["k6", "run", "/src/main.js", "--compatibility-mode=base"]