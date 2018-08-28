.PHONY: install build lib bin

default: build

LIBDIR = /usr/local/lib/geodelib

# tell go to install to the current directory
export GOBIN=$(shell pwd)

install: install.lib install.bin


clean: gc.clean
	@geode clean

build: gc.all
	@go build -o bin/geode ./pkg/cmd/geode

install.lib: gc.all
	@rm -rf $(LIBDIR)
	@mkdir -p $(LIBDIR)
	@cp -a lib/* $(LIBDIR)
	
install.bin:
	@install ./bin/geode /usr/local/bin

gen:
	@go generate -v ./...


all: build install.lib install.bin

gc.clean:
	find gc -name "*.o" | xargs rm -rf "{}"
	rm -rf lib/gc/*
	rm -rf lib/include/gc

gc.all: gc lib/gc/gc.a

lib/gc/gc.a:
	cd gc && make -f Makefile.direct
	mkdir -p lib/gc/
	cp gc/gc.a lib/gc/gc.a
	mkdir -p lib/include/gc
	cp -a gc/include/ lib/include/gc
	
	
gc:
	mkdir -p dl
	wget -O dl/libatomic_ops.tar.gz https://github.com/ivmai/libatomic_ops/releases/download/v7.6.6/libatomic_ops-7.6.6.tar.gz
	wget -O dl/gc.tar.gz https://github.com/ivmai/bdwgc/releases/download/v7.6.8/gc-7.6.8.tar.gz
	rm -rf gc
	tar -C ./ -xvf dl/gc.tar.gz
	mv gc-* gc
	tar -C ./ -xvf dl/libatomic_ops.tar.gz
	mv libatomic_ops* gc/libatomic_ops	
	rm -rf dl
