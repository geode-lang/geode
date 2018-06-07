FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go", "build"]
CMD ["file", "./act"]
CMD ["./act", "example"]