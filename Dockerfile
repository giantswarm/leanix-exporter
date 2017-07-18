FROM debian
COPY ./leanix-exporter /
ENTRYPOINT ["/leanix-exporter"]
CMD ["run"]
