# Act language

Building from source:

Download the source:

```
git clone git@github.com:nickwanninger/act.git $GOPATH/src/nickwanninger/act
cd $GOPATH/src/nickwanninger/act
```

Install the dependencies (and setup the llvm bindings). This can take some time. This also requires `svn` to be installed because llvm is still built on subversion it seems. So go install that now.

You will also need to have [dep](https://github.com/golang/dep) installed on your system

```
make deps
```

build the app to ./actc

```
make
```
