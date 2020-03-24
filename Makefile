export GO111MODULE=on
.PHONY: build

all: deps build build-docker clean

build:
	go build ./cmd/namespace-reaper

build-docker: deps build
	docker build .

deps:
	go mod tidy

clean:
	rm namespace-reaper
