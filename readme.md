<div style="text-align:center"><img src="https://s3-us-west-2.amazonaws.com/nickwanninger/geode/masthead.png"/></div>

## Installing

Installing Geode is simple, just follow the steps below and install a few dependencies

### Dependencies

- Golang with a `$GOPATH` setup in your env
- The clang c compiler for linking binaries

### Building

Once you have the dependencies setup, building is easy:

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
