# Image URL to use all building/pushing image targets
artifactURL ?= services
TAG ?= latest


start:
	go run cmd/autoscaler.go

build:
	go build -o bin/manager cmd/autoscaler.go

docker-build:
	docker build -t ${artifactURL}:${TAG} .

docker-push:
	docker push ${artifactURL}:${TAG}