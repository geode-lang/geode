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
The first build could take anywhere from 3 to 10 minutes to compile llvm bindings

```
$ go get -u gitlab.com/nickwanninger/geode
$ cd $GOPATH/src/gitlab.com/nickwanninger/geode
$ make dep
$ make
```

This will install geode's executable binary to `$GOPATH/bin/`

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
