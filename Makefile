export GO111MODULE=on
export VERSION=0.1.0
.PHONY: build

all: deps build build-docker clean

build:
	@GOOS=linux GOARCH=amd64 go build -o bin/ ./cmd/namespace-reaper 

build-arm:
	@GOOS=linux GOARCH=arm64 go build -o bin/ ./cmd/namespace-reaper 

build-docker: deps build
	docker build . -t dwardu/namespace-reaper:$(VERSION)
	docker build . -t dwardu/namespace-reaper:latest

publish: build-docker
	docker push dwardu/namespace-reaper:$(VERSION)
	docker push dwardu/namespace-reaper:latest

deps:
	go mod tidy

clean:
	rm namespace-reaper
