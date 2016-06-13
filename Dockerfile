FROM golang:1.4

RUN go get github.com/jmartin82/mmock

RUN ln -s /go/src/github.com/jmartin82/mmock/tmpl /go/bin/tmpl
RUN mkdir /config
VOLUME /config

EXPOSE 8082 8083

ENTRYPOINT ["/go/bin/mmock","-config-path","/config"]