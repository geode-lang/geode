.PHONY: install build build.bin clean docker.build docker.run docker.push




BUILDPATH = ./bin/geode
BINPATH   = /usr/local/bin/geode
# GCPATH    = lib/gc/gc.a
LIBDIR    = /usr/local/lib/geodelib

# tell go to install to the current directory
export GOBIN=$(shell pwd)

install.bin:
	@mkdir -p bin
	@install ./bin/geode /usr/local/bin

clean: lib.clean
	@rm -rf ./bin
	@geode clean

build: build.bin $(GCPATH)


build.bin:
	@go build -o ./bin/geode ./pkg/cmd/geode

uninstall: clean
	rm -rf $(shell which geode)
	rm -rf $(LIBDIR)

gen:
	@go generate -v ./...

default:

lib.clean:
	# @rm -rf gc
	@rm -rf lib/gc/*
	@rm -rf lib/include/gc

# $(GCPATH):
# 	mkdir -p dl
# 	wget -O dl/libatomic_ops.tar.gz https://github.com/ivmai/libatomic_ops/releases/download/v7.6.6/libatomic_ops-7.6.6.tar.gz
# 	wget -O dl/gc.tar.gz https://github.com/ivmai/bdwgc/releases/download/v7.6.8/gc-7.6.8.tar.gz
# 	rm -rf gc
# 	tar -C ./ -xvf dl/gc.tar.gz
# 	mv gc-* gc
# 	tar -C ./ -xvf dl/libatomic_ops.tar.gz
# 	mv libatomic_ops* gc/libatomic_ops
# 	rm -rf dl
# 	cd gc && make -f Makefile.direct
# 	mkdir -p lib/gc/
# 	cp gc/gc.a lib/gc/gc.a
# 	cd lib/gc && ar x gc.a
# 	mkdir -p lib/include/gc
# 	cp -a gc/include/* lib/include/gc
# 	rm -rf gc


install.lib:# $(GCPATH)
	@rm -rf $(LIBDIR)
	@mkdir -p $(LIBDIR)
	@cp -a lib/* $(LIBDIR)



docker.build:
	docker build -t nickwanninger/geode-test:latest .

docker.run: docker.build
	docker run --rm -it nickwanninger/geode-test:latest

docker.push: docker.build
	docker push nickwanninger/geode-test:latest


default: build
install: install.lib install.bin
all: build install