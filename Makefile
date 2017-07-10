.PHONY: all clean build run debug

all: clean build

clean:
	go clean -i ./...

run:
	go run imap-benchmark.go -logtostderr=true -v=2

debug:
	go run imap-benchmark.go -logtostderr=true -v=3

build:
	go build imap-benchmark.go
