<div style="text-align:center"><img src="https://s3-us-west-2.amazonaws.com/nickwanninger/geode/masthead.png"/></div>

[![CircleCI](https://circleci.com/gh/nickwanninger/geode/tree/master.svg?style=svg)](https://circleci.com/gh/nickwanninger/geode/tree/master)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgeode-lang%2Fgeode.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fgeode-lang%2Fgeode?ref=badge_shield)

## About

Geode is a programming language written in go around the llvm framework.
It's semi low level for the time being with plans to be higher level in
the future. It is just a compiler to llvm, then it calls clang to link the
.ll files to a runnable binary. Clang will also link the c standard library.

## Installing

Installing Geode is simple, just follow the steps below and install a few dependencies

Go to [releases](https://github.com/geode-lang/geode/releases) and download the distribution for your OS

#### MacOS

Install dependencies

```
$ brew install clang libgc
```

Run the pkg installer

#### Ubuntu/Debian

```
$ sudo apt install clang libgc-dev
$ sudo apt install ./geode-X.X.X.deb
```

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fgeode-lang%2Fgeode.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fgeode-lang%2Fgeode?ref=badge_large)

```

```

```

```
