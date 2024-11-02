# Makefile

PROTOC ?= protoc
PROTO_DIR = proto
GEN_GO_DIR = gen/go
GEN_SWAGGER_DIR = gen/swagger
THIRD_PARTY_DIR = third_party/googleapis

.PHONY: all proto clean docker-build

all: proto

proto:
	@echo "Generating Go code and Swagger documentation..."
	mkdir -p $(GEN_GO_DIR)
	mkdir -p $(GEN_SWAGGER_DIR)
	$(PROTOC) -I . \
		-I $(THIRD_PARTY_DIR) \
		--go_out=$(GEN_GO_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_GO_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GEN_GO_DIR) --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=$(GEN_SWAGGER_DIR) --openapiv2_opt=logtostderr=true \
		$(PROTO_DIR)/server.proto

clean:
	@echo "Cleaning generated code..."
	rm -rf $(GEN_GO_DIR)/*
	rm -rf $(GEN_SWAGGER_DIR)/*

docker-build:
	@echo "Building Docker images..."
	docker build -t telegram-grpc-service-server -f build/server/Dockerfile .
	docker build -t telegram-grpc-service-worker -f build/worker/Dockerfile .

proto:
	@echo "Generating code..."
	protoc -I . \
	  -I ./third_party/googleapis \
	  -I ./third_party \
	  --go_out=gen/go --go_opt=paths=source_relative \
	  --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
	  --grpc-gateway_out=gen/go --grpc-gateway_opt=paths=source_relative \
	  --openapiv2_out=gen/swagger --openapiv2_opt=logtostderr=true \
	  proto/messaging.proto

docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up
