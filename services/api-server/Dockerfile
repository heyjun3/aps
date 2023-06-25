FROM golang:1.20.5-bullseye AS builder

ARG APP_PATH
WORKDIR /app

COPY ${APP_PATH}/ /app/
RUN go build -o server

FROM gcr.io/distroless/base-debian11:latest
WORKDIR /app
COPY --from=builder /app/server /app/server
