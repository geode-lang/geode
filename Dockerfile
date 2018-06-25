FROM ubuntu
RUN mkdir -p /go
ENV GOPATH /go
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y golang-go clang