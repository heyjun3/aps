FROM golang:1.20.5-bullseye AS builder

ARG APP_PATH
WORKDIR /app

COPY ${APP_PATH}/ /app/
RUN go build -o server

FROM gcr.io/distroless/base-debian11:latest
WORKDIR /app
COPY --chown=notroot:notroot --from=builder /app/server /app/server
USER notroot
ENTRYPOINT [ "/app/server" ]
