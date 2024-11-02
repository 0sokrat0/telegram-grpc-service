# Указываем версию buf
BUF_VERSION=1.6.0

# Указываем пути
PROTO_DIR=./proto
GEN_DIR=./gen
GOPATH_BIN=$(shell go env GOPATH)/bin

.PHONY: all
all: clean format gen lint gen-go gen-openapi

# Цель для установки buf
.PHONY: buf-install
buf-install:
	@if [ ! -f $(GOPATH_BIN)/buf ]; then \
		echo "Устанавливаю buf..."; \
		tmp=$$(mktemp -d); cd $$tmp; \
		OS=$$(uname -s); \
		ARCH=$$(uname -m); \
		if [ "$$OS" = "Darwin" ]; then \
			OS="Darwin"; \
		else \
			OS="Linux"; \
		fi; \
		if [ "$$ARCH" = "x86_64" ]; then \
			ARCH="x86_64"; \
		elif [ "$$ARCH" = "arm64" ] || [ "$$ARCH" = "aarch64" ]; then \
			ARCH="arm64"; \
		else \
			echo "Unsupported architecture: $$ARCH"; exit 1; \
		fi; \
		curl -L -o buf \
			https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$$OS-$$ARCH; \
		chmod +x buf; \
		mv buf $(GOPATH_BIN)/buf; \
		echo "buf установлен в $(GOPATH_BIN)/buf"; \
	else \
		echo "buf уже установлен"; \
	fi

# Цель для генерации Go файлов из .proto
.PHONY: gen-go
gen-go:
	protoc -I . \
		-I ./third_party \
		--go_out=$(GEN_DIR)/go \
		--go-grpc_out=$(GEN_DIR)/go \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# Цель для генерации OpenAPI спецификации
.PHONY: gen-openapi
gen-openapi:
	protoc -I . \
		-I ./third_party \
		--openapiv2_out=$(GEN_DIR)/openapi \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt json_names_for_fields=false \
		$(PROTO_DIR)/*.proto

# Цель для генерации кода с помощью buf
.PHONY: gen
gen: buf-install
	buf generate $(PROTO_DIR)
	@for dir in $(GEN_DIR)/go/*/; do \
	  cd $$dir; \
	  go mod init $$(basename $$dir); \
	  go mod tidy; \
  	done

# Цель для линтинга
.PHONY: lint
lint: buf-install
	buf lint $(PROTO_DIR)

# Цель для форматирования
.PHONY: format
format: buf-install
	buf format -w $(PROTO_DIR)

# Цель для очистки
.PHONY: clean
clean:
	rm -rf $(GEN_DIR)

# Дополнительная цель для запуска сервера
.PHONY: run-server
run-server:
	go run ./cmd/server/main.go
