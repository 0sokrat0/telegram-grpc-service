package server

import (
	"github.com/0sokrat0/telegram-grpc-service/internal/api"
	"log"
	"net"

	proto_tg_service "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"google.golang.org/grpc"
)

func MainServer() {
	// Создаем новый сервис
	messagingService, err := api.NewMessagingService()
	if err != nil {
		log.Fatalf("Ошибка при создании сервиса: %v", err)
	}
	defer messagingService.Close()

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервис
	proto_tg_service.RegisterMessagingServiceServer(grpcServer, messagingService)

	// Запускаем сервер на порту 50051
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось запустить слушатель: %v", err)
	}
	log.Println("gRPC сервер запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
