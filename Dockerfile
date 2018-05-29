FROM golang:1.9-alpine as builder

# Install ca-certificates, required for the "release message" feature:
RUN apk --no-cache add \
    ca-certificates

RUN apk --no-cache add --virtual build-dependencies \
    git \
  && mkdir -p /root/gocode \
  && export GOPATH=/root/gocode \
  && go get github.com/jmartin82/mmock \
  && mv /root/gocode/bin/mmock /usr/local/bin \
  && rm -rf /root/gocode \
  && apk del --purge build-dependencies

FROM alpine:3.6

RUN apk --no-cache add \
    ca-certificates

COPY --from=builder /usr/local/bin/mmock /usr/local/bin/mmock

RUN mkdir /config
RUN mkdir /tls

VOLUME /config

ADD ./tls/server.crt /tls
ADD ./tls/server.key /tls

EXPOSE 8082 8083 8084

ENTRYPOINT ["mmock","-config-path","/config","-tls-path","/tls"]
CMD ["-server-ip","0.0.0.0","-console-ip","0.0.0.0"]
