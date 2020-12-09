.DEFAULT_GOAL := docker-image

IMAGE ?= docker.io/rdtigera/init-container:rene

image/init-container: $(shell find . -name '*.go')
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o $@ ./init-container

.PHONY: docker-image
docker-image: image/init-container
	docker build -t $(IMAGE) ./

.PHONY: push-image
push-image: docker-image
	docker push $(IMAGE)

