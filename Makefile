all: build

build: clean
	go build act.go

uninstall:
	rm -f $(GOPATH)/bin/act


install: uninstall
	go install github.com/nickwanninger/act


example: build
	./act example/example.act


clean:
	rm -rf act