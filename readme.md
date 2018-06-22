<div style="text-align:center"><img src="https://s3-us-west-2.amazonaws.com/nickwanninger/geode/masthead.png"/></div>

## Installing

Installing Geode is simple, just follow the steps below and install a few dependencies

### Dependencies

- Golang with a `$GOPATH` setup in your env
- [Go Dep](https://github.com/golang/dep) installed
- The CC compiler for linking binaries
- For building LLVM bindings:
  - Subversion
  - A C++ Compiler
  - CMake installed
  - `libedit-dev` installed

### Building

Once you have the dependencies setup, building is easy:

First you need to build the llvm bindings. This can take some time depending on the speed of your machine. Luckly this doesn't need to be done on every build.

```
$ git clone https://github.com/go-llvm/llvm.git $GOPATH/src/github.com/go-llvm/llvm
$ cd $GOPATH/src/github.com/go-llvm/llvm
$ ./update_llvm.sh
```

Then you can get the Geode Compiler

```
$ go get -u gitlab.com/nickwanninger/geode/...
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
