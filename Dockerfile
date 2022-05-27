FROM golang:alpine AS build
RUN apk add build-base

WORKDIR /src
COPY . /src

RUN go build -v \
    -ldflags="-X main.Version=$(git describe --tags) -X main.Commit=$(git rev-parse --short HEAD) -X main.Timestamp=$(date +'%Y-%m-%dT%H:%M:%S')" \
    -o /app/skeleton-service \
    ./cmd/skeleton-service/main.go;

RUN go build -v \
    -o /app/migrate \
    ./cmd/migrate/main.go;

FROM alpine
RUN apk add netcat-openbsd curl openssl bash

EXPOSE 8080

COPY --from=build /app/skeleton-service /app/skeleton-service
COPY --from=build /app/migrate /app/migrate
COPY ./cmd/skeleton-service/entrypoint.sh /app
COPY ./config/. /app/config/
COPY ./db /app/config/db

RUN adduser -s /bin/bash -D skeleton & chown -R skeleton:skeleton /app
USER skeleton
WORKDIR /app

CMD ["/app/entrypoint.sh"]