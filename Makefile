
all: build

build: clean
	@echo "Building"
	@go install  -gcflags '-N -l' gitlab.com/nickwanninger/geode/...

uninstall:
	rm -f $(GOPATH)/bin/geode

install: uninstall
	go install gitlab.com/nickwanninger/geodec
	
watch:
	nodemon --watch pkg/ --watch pkg/cmd/geode/main.go --ext go --exec make

clean:
	@rm -rf build
	@rm -rf geode

gen:
	go generate ./...

deps:
	dep ensure
	vendor/github.com/go-llvm/llvm/update_llvm.sh
