export CGO_ENABLED=0
export GOPATH?=$(shell go env GOPATH)
export DESTDIR?=$(GOPATH)/bin
export GOBIN?=$(DESTDIR)

all: build
ci: env test

env:
	go env
	echo "---"

dep:
	go get -u ./

build:
	go build
	go install github.com/xor-gate/debpkg/cmd/debpkg

test:
	go test -v $(shell go list ./... | grep -v '^vendor\/')

lint:
	go get -u github.com/golang/lint/golint
	golint ./... | grep -v '^vendor\/' | grep -v ".pb.*.go:" || true

clean:
	rm -Rf $(TMPDIR)/debpkg*

fmt:
	gofmt -s -w .
