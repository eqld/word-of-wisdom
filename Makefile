.PHONY: create-network build-server build-client run-server run-client build-and-run-server build-and-run-client stop-server stop-client remove-network

SERVER_IMAGE_NAME = word-of-wisdom-server
CLIENT_IMAGE_NAME = word-of-wisdom-client
SERVER_CONTAINER_NAME = wow-server
CLIENT_CONTAINER_NAME = wow-client
NETWORK_NAME = wow-network
SERVER_PORT = 8080
CLIENT_HOST ?= $(SERVER_CONTAINER_NAME)
CLIENT_PORT ?= 8080
WOW_SERVER_DIFFICULTY ?= 2
WOW_SERVER_CHALLENGE_LENGTH ?= 16
WOW_SERVER_SOLUTION_LENGTH ?= 8
WOW_CLIENT_SOLUTION_LENGTH ?= $(WOW_SERVER_SOLUTION_LENGTH)

create-network:
	-docker network create $(NETWORK_NAME)

build-server:
	docker build -t $(SERVER_IMAGE_NAME) -f server.Dockerfile .

build-client:
	docker build -t $(CLIENT_IMAGE_NAME) -f client.Dockerfile .

run-server: create-network
	docker run -it --rm --name $(SERVER_CONTAINER_NAME) \
		--env WOW_SERVER_DIFFICULTY=$(WOW_SERVER_DIFFICULTY) \
		--env WOW_SERVER_CHALLENGE_LENGTH=$(WOW_SERVER_CHALLENGE_LENGTH) \
		--env WOW_SERVER_SOLUTION_LENGTH=$(WOW_SERVER_SOLUTION_LENGTH) \
		--network $(NETWORK_NAME) \
		-p $(SERVER_PORT):$(SERVER_PORT) \
		$(SERVER_IMAGE_NAME)

run-client: create-network
	docker run -it --rm --name $(CLIENT_CONTAINER_NAME) \
		--env WOW_CLIENT_SOLUTION_LENGTH=$(WOW_CLIENT_SOLUTION_LENGTH) \
		--network $(NETWORK_NAME) \
		$(CLIENT_IMAGE_NAME) "$(CLIENT_HOST):$(CLIENT_PORT)"

build-and-run-server: build-server run-server

build-and-run-client: build-client run-client

stop-server:
	docker stop $(SERVER_CONTAINER_NAME)
	-docker rm $(SERVER_CONTAINER_NAME)

stop-client:
	docker stop $(CLIENT_CONTAINER_NAME)
	-docker rm $(CLIENT_CONTAINER_NAME)

remove-network:
	-docker network rm $(NETWORK_NAME)
