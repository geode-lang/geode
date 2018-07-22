
all: build

LIBDIR = /usr/local/lib/geodelib
export GOBIN=/usr/local/bin
build:
	@rm -rf $(LIBDIR)
	@mkdir -p $(LIBDIR)
	@cp -a lib/* $(LIBDIR)
	@go install github.com/nickwanninger/geode/...

install: gen build

gen:
	go generate -v ./...

dev: build
	@geode run -S example