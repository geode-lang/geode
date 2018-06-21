FLAGS := -ldflags "-X main.revhash=`git rev-parse HEAD`"


all: build

build: clean
	@go build -v -o actc main.go

uninstall:
	rm -f $(GOPATH)/bin/actc

install: uninstall
	go install gitlab.com/nickwanninger/geodec
	
watch:
	nodemon --watch pkg/ --watch main.go --ext go --exec make

example: build
	./actc example

clean:
	@rm -rf build
	@rm -rf actc


deps:
	dep ensure
	vendor/github.com/go-llvm/llvm/update_llvm.sh
