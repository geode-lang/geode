FLAGS := -ldflags "-X main.revhash=`git rev-parse HEAD`"


all: build

build: clean
	go build $(FLAGS) -o actc main.go

uninstall:
	rm -f $(GOPATH)/bin/actc


install: uninstall
	go install github.com/nickwanninger/actc


example: build
	./actc example


clean:
	rm -rf build
	rm -rf actc


# Build for every system (64 bit)
PLATFORMS := linux windows darwin plan9 netbsd openbsd freebsd
$(PLATFORMS):
	@mkdir -p build
	GOOS=$@ GOARCH=386 go build $(FLAGS) -o 'build/act-$@'

release: $(PLATFORMS)
