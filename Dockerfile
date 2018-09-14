FROM alpine:3.8

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ADD ./leanix-exporter /leanix-exporter

ENTRYPOINT ["/leanix-exporter"]

CMD [ "daemon --server.listen.address=http://0.0.0.0:8000" ]
