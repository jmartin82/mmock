FROM alpine:3.20

RUN apk --no-cache add \
    ca-certificates curl

RUN mkdir /config
RUN mkdir /tls

VOLUME /config

ARG TARGETPLATFORM
COPY tls/server.crt /tls/server.crt
COPY tls/server.key /tls/server.key
COPY ${TARGETPLATFORM}/mmock /usr/local/bin/mmock

EXPOSE 8082 8083 8084

ENTRYPOINT ["mmock","-config-path","/config","-tls-path","/tls"]
CMD ["-server-ip","0.0.0.0","-console-ip","0.0.0.0"]
HEALTHCHECK --interval=30s --timeout=3s --start-period=3s --retries=2 CMD curl -fsS http://localhost:8082 || exit 1
