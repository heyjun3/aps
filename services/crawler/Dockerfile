FROM golang:1.19.4 AS builder

ARG APP_PATH
WORKDIR /app

COPY ${APP_PATH}/ /app/
RUN go build -o go-crawler

FROM gcr.io/distroless/base-debian11:latest
WORKDIR /app
COPY --from=builder /app/go-crawler /app/go-crawler
COPY --from=builder /app/sqlboiler.toml /app/sqlboiler.toml
