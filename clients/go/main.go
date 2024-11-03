package main

import (
	"context"
	"log"
	"time"

	proto "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"google.golang.org/grpc"
)

func main() {
	// Устанавливаем соединение с gRPC сервером
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := proto.NewMessagingServiceClient(conn)

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Пример: запрос на отправку текстового сообщения всем пользователям
	resp, err := client.SendMessage(ctx, &proto.SendMessageRequest{
		All: true,
		Content: &proto.SendMessageRequest_TextContent{
			TextContent: &proto.TextContent{
				Text:                  "Привет от бота! Это тестовое сообщение.",
				ParseMode:             "HTML",
				DisableWebPagePreview: true,
			},
		},
	})

	if err != nil {
		log.Fatalf("Ошибка при вызове SendMessage: %v", err)
	}
	log.Printf("Ответ от сервера: %v", resp)
}
