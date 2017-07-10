FROM golang:latest

RUN mkdir -p /go/src/github.com/go-pluto/benchmark
RUN mkdir /config
RUN mkdir /users

ADD . /go/src/github.com/go-pluto/benchmark/

WORKDIR /go/src/github.com/go-pluto/benchmark
RUN go get ./...
RUN go build imap-benchmark.go

ENTRYPOINT ["/go/src/github.com/go-pluto/benchmark"]
CMD ["-config", "/config/config.toml","-userdb","/users/passwd","-logtostderr=true","-v=3"]
