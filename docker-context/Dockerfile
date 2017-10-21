FROM golang:1.9.1-alpine3.6

LABEL Maintainer="Krzysztof Romanowski <krzysztof.romanowski94@gmail.com>"

RUN apk update && \
    apk add git && \
    apk add ansible

WORKDIR /root

VOLUME ["/root/vars"]

EXPOSE 13334

COPY ["startup.yml", "startup.yml"]
COPY ["entrypoint.sh", "/usr/bin/entrypoint"]

ENTRYPOINT ["/usr/bin/entrypoint"]

