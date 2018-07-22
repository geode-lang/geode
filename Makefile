.PHONY: install build

default: build

BINDIR =
LIBDIR = /usr/local/lib/geodelib

# tell go to install to the current directory
export GOBIN=$(shell pwd)

install:
	rm -rf $(LIBDIR)
	mkdir -p $(LIBDIR)
	cp -a lib/* $(LIBDIR)
	install ./bin/geode /usr/local/bin

build:
	go build -o bin/geode ./pkg/cmd/geode

gen:
	go generate -v ./...

dev: lib install