package api

import (
	"context"
	"fmt"
	"log"

	proto_tg_service "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// MessagingService реализует интерфейс вашего gRPC сервиса
type MessagingService struct {
	proto_tg_service.UnimplementedMessagingServiceServer
	nc *nats.Conn
	js nats.JetStreamContext
}

// NewMessagingService создает новый экземпляр MessagingService
func NewMessagingService() (*MessagingService, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("ошибка инициализации JetStream: %v", err)
	}

	return &MessagingService{
		nc: nc,
		js: js,
	}, nil
}

// Close закрывает соединение с NATS
func (m *MessagingService) Close() {
	m.nc.Close()
}

// SendMessage принимает запрос от клиента
func (m *MessagingService) SendMessage(ctx context.Context, req *proto_tg_service.SendMessageRequest) (*proto_tg_service.SendMessageResponse, error) {
	// Валидация входных данных
	if len(req.UserIds) == 0 {
		return nil, fmt.Errorf("список UserIds не должен быть пустым")
	}
	if req.GetTextContent() == nil && req.GetPhotoContent() == nil {
		return nil, fmt.Errorf("необходимо указать TextContent или PhotoContent")
	}

	// Серилизуем запрос в байты
	msgData, err := proto.Marshal(req)
	if err != nil {
		log.Printf("Ошибка сериализации запроса: %v", err)
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}

	// Публикуем сообщение в поток "MESSAGES"
	pubAck, err := m.js.Publish("MESSAGES.send_message", msgData)
	if err != nil {
		log.Printf("Ошибка публикации в JetStream: %v", err)
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}

	if pubAck == nil || pubAck.Sequence == 0 {
		log.Printf("Не удалось получить подтверждение публикации")
		return nil, fmt.Errorf("внутренняя ошибка сервера")
	}

	// Возвращаем немедленный ответ клиенту
	return &proto_tg_service.SendMessageResponse{
		Success: true,
	}, nil
}
