
all: build

build: clean
	@go install github.com/nickwanninger/geode/...

install: gen build
	
watch:
	nodemon --watch pkg/ --watch pkg/cmd/geode/main.go --ext go --exec make

clean:
	@rm -rf build
	@rm -rf geode
	@rm -rf *.s *.ll

gen:
	go generate -v ./...
