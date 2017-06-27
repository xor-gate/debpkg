export GOPATH?=$(PWD)/../../../../
export DESTDIR?=$(GOPATH)/bin
export GOBIN?=$(DESTDIR)

all: build
ci: test

build:
	go build
	go install github.com/xor-gate/debpkg/cmd/debpkg

test:
	go test -v

lint:
	go tool vet .

fmt:
	gofmt -s -w .

clean:
	rm -Rf *.deb
	rm -Rf *.tar.gz

.PHONY: clean
