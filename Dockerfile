FROM golang:1.7

RUN curl https://glide.sh/get | sh 
RUN go get github.com/jmartin82/mmock
RUN glide install
RUN mkdir /config

VOLUME /config
RUN mkdir /data
VOLUME /data

EXPOSE 8082 8083

ENTRYPOINT ["/go/bin/mmock","-config-path","/config"]