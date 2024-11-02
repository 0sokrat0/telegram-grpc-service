package main

import (
	"fmt"
	"telegram-grpc-service/config" // замените `myproject` на фактическое имя вашего модуля
)

func main() {
	cfg := config.GetConfig()

	// Используйте токен бота
	fmt.Printf("Telegram Bot Token: %s\n", cfg.BotToken)

	// Используйте параметры для gRPC
	grpcAddress := fmt.Sprintf("%s:%d", cfg.GrpcHost, cfg.GrpcPort)
	fmt.Printf("gRPC Server Address: %s\n", grpcAddress)

	// Здесь можно продолжить настройку gRPC-сервера с использованием cfg.GrpcHost и cfg.GrpcPort
}
