<div style="text-align:center"><img src="https://s3-us-west-2.amazonaws.com/nickwanninger/geode/masthead.png"/></div>

[![CircleCI](https://circleci.com/gh/nickwanninger/geode/tree/master.svg?style=svg)](https://circleci.com/gh/nickwanninger/geode/tree/master)

## About

Geode is a programming language written in go around the llvm framework.
It's semi low level for the time being with plans to be higher level in
the future. It is just a compiler to llvm, then it calls clang to link the
.ll files to a runnable binary. Clang will also link the c standard library.

You can download (semi-regularly updated) binaries from the releases section,
but you might want to just install from source regardless. This is because
the compiler relies on the library files being in the `lib/` folder inside
the $GOPATH.

Geode is a heavy work in progress with apis that will change. Extended use is
not recommended at this stage.

## Installing

Installing Geode is simple, just follow the steps below and install a few dependencies

### Dependencies

- Golang with a `$GOPATH` setup in your env
- The clang c compiler for linking binaries

### Building

Once you have the dependencies setup, building is easy:

```
$ go get -u -v github.com/nickwanninger/geode/...
```

This will build and install geode's executable binary to `$GOPATH/bin/`

## Example usage:

Geode is a massive work in progress, but you can look at example/main.g for a working state program

### Compiling a program

```
$ geode build <sourcefile>
```

Files can be any of the following:

- A folder containing a main.g
- A geode source file without the .g extension
- A .g file

### helloworld.g

```go
include "std::io"

func main int {
	print("Hello, world\n");
	return 0;
}
```

### Example fib.g

```go
include "std::io"

func fib(int n) int {
	if n < 2 {
		return n;
	}
	return fib(n - 1) + fib(n - 2);
}

func main int {
	print("%d\n", fib(30));
	return 0;
}
```
