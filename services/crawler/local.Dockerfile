FROM golang:1.21.1-bookworm

RUN apt-get update -y && \
     apt-get install -y postgresql-client iputils-ping
