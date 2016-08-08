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

test-dpkg:
	dpkg --info debpkg-test.deb
	dpkg --contents debpkg-test.deb
	dpkg --info debpkg-test-add-directory.deb
	dpkg --contents debpkg-test-add-directory.deb
	dpkg --info debpkg-test-signed.deb
	dpkg --contents debpkg-test-signed.deb

lint:
	go tool vet .

fmt:
	gofmt -s -w .

clean:
	rm -Rf *.deb

.PHONY: clean
