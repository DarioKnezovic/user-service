# Docker image name
IMAGE_NAME := user-service

.PHONY: build
build:
	docker-compose build

.PHONY: run
run:
	docker run -it --rm --name $(IMAGE_NAME) -p 5050:5050 $(IMAGE_NAME)

.PHONY: hot-reload
hot-reload:
	docker run -it --rm --name $(IMAGE_NAME) -v $(PWD):/app -w /app -p 5050:5050 $(IMAGE_NAME) air

.PHONY: dev
dev:
	docker-compose up

proto-generate:
	cd proto && protoc --go_out=. --go-grpc_out=. user.proto