FROM circleci/golang:1.10
FROM rsmmr/clang
# RUN mkdir -p /go
# ENV GOPATH /go

CMD ["/bin/ls", "/bin"]

