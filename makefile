.PHONY: build doc fmt lint dev test vet godep install bench bindata

PKG_NAME=$(shell basename `pwd`)
NS = jordimartin
VERSION ?= latest

install:
	go get -t -v ./...

build: bindata \
	vet \
	test
	go build -v -o ./bin/$(PKG_NAME)

bindata:
	go-bindata -pkg console -o console/bindata.go tmpl/*

doc:
	godoc -http=:6060

fmt:
	go fmt ./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./...

test:
	go test -v ./...
	
coverage:
	echo 'mode: atomic' > coverage.txt && go list ./... | xargs -n1 -I{} sh -c 'go test -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt'
	rm coverage.tmp

# https://godoc.org/golang.org/x/tools/cmd/vet
vet:
	go vet -v

get-deps:
	go get github.com/mattn/goveralls
	go get -u github.com/AlekSi/gocoverutil
	glide install

release:
	docker build --no-cache=true  -t $(NS)/$(PKG_NAME):$(VERSION) .
	docker push $(NS)/$(PKG_NAME):$(VERSION)

release-beta:
	docker build --no-cache=true  -t $(NS)/$(PKG_NAME):beta .
	docker push $(NS)/$(PKG_NAME):beta
