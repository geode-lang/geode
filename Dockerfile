FROM rsmmr/clang as CLANG
FROM circleci/golang:1.10

COPY --from=CLANG /opt/clang/bin/clang /opt/clang/bin/clang


ENV PATH /opt/clang/bin:$PATH

# RUN mkdir -p /go
# ENV GOPATH /go

CMD ["/bin/bash"]

