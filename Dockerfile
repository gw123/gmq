FROM alpine:3.10

COPY dist/gateway /usr/local/bin/gateway
COPY resources/public /usr/local/bin/resources/public
COPY resources/views /usr/local/bin/resources/views
COPY config.docker.yml /usr/local/bin/config/config.yml
COPY docker-entrypoint /usr/local/bin/docker-entrypoint
WORKDIR /usr/local/bin
CMD ["/usr/local/bin/gateway" ,"-c" ,"/usr/local/bin/config/config.yml"]