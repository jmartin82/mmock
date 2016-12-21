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
	go-ls -exec="go fmt" ./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	go-ls -exec="golint" ./...

test:
	go-ls -exec="go test -v" ./...
	
coverage:
	goverage -v -covermode count -coverprofile=coverage.out
	go tool cover -html=coverage.out


# https://godoc.org/golang.org/x/tools/cmd/vet
vet:
	go-ls -exec="go vet -v"

get-deps:
	go get github.com/laher/gols/cmd/go-ls
	glide install

release:
	docker build -t $(NS)/$(PKG_NAME):$(VERSION) .
	docker push $(NS)/$(PKG_NAME):$(VERSION)