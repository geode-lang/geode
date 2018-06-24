<div style="text-align:center"><img src="https://s3-us-west-2.amazonaws.com/nickwanninger/geode/masthead.png"/></div>

## Installing

Installing Geode is simple, just follow the steps below and install a few dependencies

### Dependencies

- Golang with a `$GOPATH` setup in your env
- The clang c compiler for linking binaries

### Building

Once you have the dependencies setup, building is easy:

```
$ go get -u -v gitlab.com/nickwanninger/geode/...
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



### Example fib.g

```go
func fib(int n) int {
	if n < 2 {
		return n;
	}
	return fib(n - 1) + fib(n - 2);
}

func main(int argc) byte {
	printf("%d\n", fib(30));
	return 0;
}
```