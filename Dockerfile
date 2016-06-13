FROM golang:1.4

RUN go get github.com/jmartin82/mmock

EXPOSE 8082 8083

ENTRYPOINT ["/go/bin/mmock"]