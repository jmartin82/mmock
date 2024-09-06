.PHONY: build doc fmt lint dev test vet

PKG_NAME=mmock
NS = jordimartin
VERSION ?= latest

export GO111MODULE=on


build: vet \
	test
	go build  -v -o ./bin/$(PKG_NAME) cmd/mmock/main.go

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
	go vet -v  ./...

release:
	goreleaser --clean

docker-push:
	docker build --no-cache=true  -t $(NS)/$(PKG_NAME):$(VERSION) .
	docker push $(NS)/$(PKG_NAME):$(VERSION)
