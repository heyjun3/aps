FROM golang:1.22.4-bookworm

RUN apt-get update -y && \
     apt-get install -y postgresql-client iputils-ping
RUN wget -O - https://github.com/sqldef/sqldef/releases/latest/download/psqldef_linux_amd64.tar.gz \
     | tar xvz && \
     chmod +x psqldef && \
     mv ./psqldef /usr/local/bin/
