SHELL := /bin/bash
BINARY_NAME=signer
DOCKER_COMPOSE_FILE=docker-compose.yml
DOCKER_COMPOSE=docker-compose -f $(DOCKER_COMPOSE_FILE)

.PHONY: proto-gen
proto-gen:
	docker build -t proto-genc -f ./proto/Dockerfile .
	docker run -v $(shell pwd):/src proto-genc /bin/sh scripts/proto-gen.sh

.PHONY: build
build:
	go build -o ./build/${BINARY_NAME} github.com/skip-mev/platform-take-home/cmd/signer

.PHONY: start
start:
	./build/${BINARY_NAME}

.PHONY: test
test:
	go test -v -p 1 -count 1 -race ./...

.PHONY: tf-init
tf-init:
	cd contrib/terraform && terraform init

.PHONY: tf-apply
tf-apply:
	cd contrib/terraform && terraform apply

.PHONY: start-vault
start-vault:
	vault server -dev
