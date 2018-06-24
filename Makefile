
all: build

build: clean
	@echo "Building"
	go install -v -gcflags="-N -l" gitlab.com/nickwanninger/geode/...

install: gen build
	
watch:
	nodemon --watch pkg/ --watch pkg/cmd/geode/main.go --ext go --exec make

clean:
	@rm -rf build
	@rm -rf geode

gen:
	go generate -v ./...
