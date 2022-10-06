FROM fluent/fluentd:v1.15-debian-1

USER root
RUN apt update -y && apt install -y make gcc
RUN gem install fluent-plugin-mongo
COPY ./docker/fluentd/http.conf /fluent/etc/

CMD [ "fluentd", "-c", "/fluent/etc/http.conf" ]
