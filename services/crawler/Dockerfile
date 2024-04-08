FROM golang:1.21.1-bookworm AS builder

ARG APP_PATH
WORKDIR /app

COPY ${APP_PATH}/ /app/
RUN go build -o go-crawler
RUN go build -o serve cmd/server/main.go
RUN go build -o migrate db/cmd/main.go

FROM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --from=builder /app/go-crawler /app/go-crawler
COPY --from=builder /app/serve /app/server
COPY --from=builder /app/migrate /app/migrate
