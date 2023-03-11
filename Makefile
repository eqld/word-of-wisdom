.PHONY: create-network build-server build-client run-server run-client build-and-run-server build-and-run-client stop-server stop-client remove-network

SERVER_IMAGE_NAME = word-of-wisdom-server
CLIENT_IMAGE_NAME = word-of-wisdom-client
SERVER_CONTAINER_NAME = wow-server
CLIENT_CONTAINER_NAME = wow-client
NETWORK_NAME = wow-network
SERVER_PORT = 8080
CLIENT_HOST ?= $(SERVER_CONTAINER_NAME)
CLIENT_PORT ?= 8080

create-network:
	-docker network create $(NETWORK_NAME)

build-server:
	docker build -t $(SERVER_IMAGE_NAME) -f server.Dockerfile .

build-client:
	docker build -t $(CLIENT_IMAGE_NAME) -f client.Dockerfile .

run-server: create-network
	docker run -it --rm --name $(SERVER_CONTAINER_NAME) --network $(NETWORK_NAME) -p $(SERVER_PORT):$(SERVER_PORT) $(SERVER_IMAGE_NAME)

run-client: create-network
	docker run -it --rm --name $(CLIENT_CONTAINER_NAME) --network $(NETWORK_NAME) $(CLIENT_IMAGE_NAME) "$(CLIENT_HOST):$(CLIENT_PORT)"

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
