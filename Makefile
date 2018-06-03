all: build

build: clean
	go build -o actc main.go

uninstall:
	rm -f $(GOPATH)/bin/actc


install: uninstall
	go install github.com/nickwanninger/actc


example: build
	./actc example/example.act


clean:
	rm -rf actc