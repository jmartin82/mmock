FROM golang:1.4

RUN go get github.com/jmartin82/mmock

RUN mkdir /config
VOLUME /config

EXPOSE 8082 8083

ENTRYPOINT ["/go/bin/mmock","-config-path","/config"]