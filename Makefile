export GOPATH?=$(PWD)/../../../../
export DESTDIR?=$(GOPATH)/bin
export GOBIN?=$(DESTDIR)

all: build
ci: test

build:
	go build
	go install github.com/xor-gate/debpkg/cmd/debpkg

test:
	go test -v -race

lint:
	go tool vet .

fmt:
	gofmt -s -w .

clean:
	rm -Rf *.deb

.PHONY: clean
