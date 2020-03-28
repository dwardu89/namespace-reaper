export GO111MODULE=on
.PHONY: build

all: deps build build-docker clean

build-minikube: deps build
	eval $(minikube docker-env)
	docker build . -t namespace-reaper

deploy-minikube: build-minikube
	kubectl apply -f minikube-pod.yaml

build:
	@GOOS=linux GOARCH=amd64 go build ./cmd/namespace-reaper

build-docker: deps build
	docker build . -t dwardu/namespace-reaper

deps:
	go mod tidy

clean:
	rm namespace-reaper
	