.PHONY: install build lib bin

default: build

LIBDIR = /usr/local/lib/geodelib

# tell go to install to the current directory
export GOBIN=$(shell pwd)

install: lib bin

build:
	@go build -v -o bin/geode ./pkg/cmd/geode

lib:
	@rm -rf $(LIBDIR)
	@mkdir -p $(LIBDIR)
	@cp -a lib/* $(LIBDIR)
	
bin:
	@install ./bin/geode /usr/local/bin

gen:
	go generate -v ./...


all: build install
