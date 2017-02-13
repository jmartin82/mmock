FROM alpine:3.4

# Install ca-certificates, required for the "release message" feature:
RUN apk --no-cache add \
    ca-certificates

# Install MailHog:
RUN apk --no-cache add --virtual build-dependencies \
    go \
    git \
  && mkdir -p /root/gocode \
  && export GOPATH=/root/gocode \
  && go get github.com/jmartin82/mmock \
  && mv /root/gocode/bin/mmock /usr/local/bin \
  && rm -rf /root/gocode \
  && apk del --purge build-dependencies

RUN mkdir /config

VOLUME /config

EXPOSE 8082 8083

ENTRYPOINT ["mmock","-config-path","/config"]