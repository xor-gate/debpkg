export GOPATH?=$(PWD)/../../../../
export DESTDIR?=$(GOPATH)/bin
export GOBIN?=$(DESTDIR)

all: build
ci: lint test

build:
	go build

test:
	go test -v -race

lint:
	go tool vet .

fmt:
	gofmt -d -s .
