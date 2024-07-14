FROM golang:1.22.4-bookworm AS base

ARG APP_PATH
WORKDIR /app

COPY ${APP_PATH}/ /app/

FROM base AS build-crawler
RUN go build -o go-crawler

FROM base AS build-server
RUN go build -o serve cmd/server/main.go

FROM base AS build-migrate
RUN go build -o migrate cmd/migrate/main.go

FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=build-crawler /app/go-crawler /app/go-crawler
COPY --from=build-server /app/serve /app/server
COPY --from=build-migrate /app/migrate /app/migrate
