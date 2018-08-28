.PHONY: install build lib bin

default: build

LIBDIR = /usr/local/lib/geodelib

# tell go to install to the current directory
export GOBIN=$(shell pwd)

install: lib bin clean


clean: gc.clean
	@geode clean

build:
	@go build -o bin/geode ./pkg/cmd/geode

lib: gc
	@rm -rf $(LIBDIR)
	@mkdir -p $(LIBDIR)
	@cp -a lib/* $(LIBDIR)
	
bin:
	@install ./bin/geode /usr/local/bin

gen:
	@go generate -v ./...


all: build lib bin

gc.clean:
	find gc -name "*.o" | xargs rm -rf "{}"
	rm -rf lib/gc/*
	rm -rf lib/include/gc

gc: lib/gc/gc.a

lib/gc/gc.a:
	cd gc && make -f Makefile.direct
	cp gc/gc.a lib/gc/gc.a
	mkdir -p lib/include/gc
	cp -a gc/include/ lib/include/gc