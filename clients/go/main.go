package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	proto_tg_service "telegram-grpc-service/gen/go/proto"
)

func main() {
	// Устанавливаем соединение с gRPC сервером
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := proto_tg_service.NewMessagingServiceClient(conn)

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Список Telegram User IDs для отправки сообщения
	userIDs := []int64{575225733} // Замените на реальные User IDs

	// Отправляем текстовое сообщение
	resp, err := client.SendMessage(ctx, &proto_tg_service.SendMessageRequest{
		UserIds: userIDs,
		Content: &proto_tg_service.SendMessageRequest_TextContent{
			TextContent: &proto_tg_service.TextContent{
				Text:                  "Привет от бота!",
				ParseMode:             "MarkdownV2",
				DisableWebPagePreview: true,
			},
		},
	})
	if err != nil {
		log.Fatalf("Ошибка при вызове SendMessage: %v", err)
	}

	log.Printf("Ответ от сервера: %v", resp)
}
