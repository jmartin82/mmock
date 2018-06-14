.PHONY: build doc fmt lint dev test vet godep install bench bindata

PKG_NAME=mmock
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
	./coverage.sh

# https://godoc.org/golang.org/x/tools/cmd/vet
vet:
	go vet -v

get-deps:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

release:
	goreleaser --rm-dist
	docker build --no-cache=true  -t $(NS)/$(PKG_NAME):$(VERSION) .
	docker push $(NS)/$(PKG_NAME):$(VERSION)
