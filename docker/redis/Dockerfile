FROM redis:7.0.5

COPY ./docker/redis/redis.conf /etc/redis.conf

ENTRYPOINT [ "redis-server", "/etc/redis.conf" ]
