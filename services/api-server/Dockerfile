FROM golang:1.21.4-bullseye AS builder

ARG APP_PATH
WORKDIR /app

COPY ${APP_PATH}/ /app/
RUN go build -o server
RUN go build cmd/delete_row.go

FROM gcr.io/distroless/base-debian11:latest
WORKDIR /app
COPY --chown=nonroot:nonroot --from=builder /app/server /app/server
COPY --chown=nonroot:nonroot --from=builder /app/delete_row /app/delete_row
COPY --chown=nonroot:nonroot --from=builder /app/database/migrations /app/database/migrations
USER nonroot
