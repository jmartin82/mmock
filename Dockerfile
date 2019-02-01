FROM alpine:3.6

RUN apk --no-cache add \
    ca-certificates

RUN mkdir /config
RUN mkdir /tls

VOLUME /config

COPY tls/server.crt /tls/server.crt
COPY tls/server.key /tls/server.key
COPY mmock /usr/local/bin/mmock

EXPOSE 8082 8083 8084

ENTRYPOINT ["mmock","-config-path","/config","-tls-path","/tls"]
CMD ["-server-ip","0.0.0.0","-console-ip","0.0.0.0"]
