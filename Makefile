.PHONY: all clean build run debug

all: clean build

clean:
	go clean -i ./...

build:
	CGO_ENABLED=0 go build -ldflags '-extldflags "-static"'

run:
	go run main.go -logtostderr=true -v=2

debug:
	go run main.go -logtostderr=true -v=3
