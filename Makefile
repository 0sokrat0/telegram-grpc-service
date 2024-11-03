PROTOC ?= protoc
PROTO_DIR = proto
GEN_GO_DIR = gen/go
GEN_SWAGGER_DIR = gen/swagger
THIRD_PARTY_DIR = third_party/googleapis

.PHONY: all proto clean docker-build docker-up

all: proto

proto:
	@echo "Generating Go code and Swagger documentation..."
		 mkdir -p gen/go
		 mkdir -p gen/swagger
		 protoc -I . \
			-I third_party \
			-I third_party/googleapis \
			-I third_party/protoc-gen-openapiv2 \
			--go_out=gen/go --go_opt=paths=source_relative \
			--go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
			--grpc-gateway_out=gen/go --grpc-gateway_opt=paths=source_relative \
			--openapiv2_out=gen/swagger --openapiv2_opt=logtostderr=true \
			proto/messaging.proto

clean:
	@echo "Cleaning generated code..."
	rm -rf $(GEN_GO_DIR)/*
	rm -rf $(GEN_SWAGGER_DIR)/*

docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up
