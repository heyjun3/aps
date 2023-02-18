FROM golang:1.19.4

RUN apt-get update -y && apt-get install -y postgresql-client
