#!/bin/bash

# Остановить скрипт при ошибке
set -e

# Определите переменные
NATS_COMPOSE_FILE="docker-compose.nats.yml"
APP_BINARY="main"
APP_DIR="telegram-grpc-service"  # Замените на путь к вашему проекту
APP_ENV_FILE=".env"  # Если у вас есть .env файл с переменными окружения

# Создание Docker Compose файла для NATS
cat <<EOL > $NATS_COMPOSE_FILE
version: "3.8"

services:
  nats-server:
    image: nats:latest
    container_name: nats-server
    networks:
      - app_network
    ports:
      - "4222:4222"
      - "8222:8222"

networks:
  app_network:
    name: app_network
    driver: bridge
EOL

# Запуск NATS с помощью Docker Compose
echo "Запуск NATS..."
docker-compose -f $NATS_COMPOSE_FILE up -d

# Сборка Go-приложения
echo "Сборка Go-приложения..."
cd "$PWD"
go build -o $APP_BINARY ./cmd/main

# Установка переменных окружения, если есть файл .env
if [ -f "$APP_ENV_FILE" ]; then
  export $(grep -v '^#' $APP_ENV_FILE | xargs)
fi

# Запуск Go-приложения
echo "Запуск Go-приложения..."
./$APP_BINARY &

echo "Проект успешно развернут!"
