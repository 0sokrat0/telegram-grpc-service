package main

import (
	"context"
	proto_tg_service "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"google.golang.org/grpc"
	"log"
	"time"
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

	// Отправляем сообщение с фото и форматированной подписью
	resp, err := client.SendMessage(ctx, &proto_tg_service.SendMessageRequest{
		UserIds: userIDs,
		Content: &proto_tg_service.SendMessageRequest_PhotoContent{
			PhotoContent: &proto_tg_service.PhotoContent{
				Url:       "https://miro.medium.com/v2/resize:fit:1000/0*YISbBYJg5hkJGcQd.png", // Замените на реальный URL изображения
				Caption:   "Привет от бота! <b>Это жирный текст</b>, <i>это курсив</i>. <a href=\"https://example.com\">Ссылка</a>",
				ParseMode: "HTML",
			},
		},
	})

	if err != nil {
		log.Fatalf("Ошибка при вызове SendMessage: %v", err)
	}
	log.Printf("Ответ от сервера: %v", resp)
}
