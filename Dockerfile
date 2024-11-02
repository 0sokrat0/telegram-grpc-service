# build/gateway/Dockerfile

# Stage 1: Сборка Go-бинарника
FROM golang:1.19-alpine AS builder

WORKDIR /app

# Копируем файлы модулей и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway ./cmd/gateway

# Stage 2: Создание минимального образа
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из предыдущего этапа
COPY --from=builder /app/gateway .

# Открываем порт для HTTP сервера
EXPOSE 8080

# Запуск gateway
CMD ["./gateway", "--grpc-server-endpoint=grpc_server:50051"]
