FLAGS := -ldflags "-X main.revhash=`git rev-parse HEAD`"


all: build

build: clean
	go build -v -o actc main.go

uninstall:
	rm -f $(GOPATH)/bin/actc


install: uninstall
	go install github.com/nickwanninger/actc


example: build
	./actc example


clean:
	rm -rf build
	rm -rf actc
	
	
deps:
	dep ensure
	vendor/github.com/nickwanninger/llvm/update_llvm.sh
