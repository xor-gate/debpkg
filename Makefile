export GOPATH?=$(PWD)/../../../../
export DESTDIR?=$(GOPATH)/bin
export GOBIN?=$(DESTDIR)

all: build
ci: test

dep:
	go get -u ./

build:
	go build
	go install github.com/xor-gate/debpkg/cmd/debpkg

test:
	go test -v

lint:
	go get -u github.com/golang/lint/golint
	golint ./... | grep -v '^vendor\/' | grep -v ".pb.*.go:" || true

clean:
	rm -Rf $(TMPDIR)/debpkg*

fmt:
	gofmt -s -w .
