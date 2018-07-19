
all: build

build:
	@go install github.com/nickwanninger/geode/...

install: gen build


gen:
	go generate -v ./...

dev: build
	@geode run -S example