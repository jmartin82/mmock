FROM golang:1.4

RUN go get github.com/jmartin82/mmock

RUN mkdir /config
VOLUME /config

RUN mkdir /data
VOLUME /data

EXPOSE 8082 8083

ENTRYPOINT ["/go/bin/mmock","-config-path","/config","-config-persist-path","/data"]