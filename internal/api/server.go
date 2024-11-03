package api

import (
	"context"
	"fmt"
	proto_tg_service "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"github.com/0sokrat0/telegram-grpc-service/internal/database"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"log"
)

type MessagingService struct {
	proto_tg_service.UnimplementedMessagingServiceServer
	nc *nats.Conn
	js nats.JetStreamContext
}

// NewMessagingService создает новый экземпляр MessagingService
func NewMessagingService() (*MessagingService, error) {
	natsURL := "nats://0.0.0.0:4222"
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("ошибка инициализации JetStream: %v", err)
	}

	database.InitDB()

	return &MessagingService{
		nc: nc,
		js: js,
	}, nil
}

// Close закрывает соединение с NATS
func (m *MessagingService) Close() {
	m.nc.Close()
}

// SendMessage отправляет сообщение пользователям
func (m *MessagingService) SendMessage(ctx context.Context, req *proto_tg_service.SendMessageRequest) (*proto_tg_service.SendMessageResponse, error) {
	if !req.All && len(req.UserIds) == 0 {
		return nil, fmt.Errorf("список UserIds не должен быть пустым, если all=false")
	}
	if req.GetTextContent() == nil && req.GetPhotoContent() == nil {
		return nil, fmt.Errorf("необходимо указать TextContent или PhotoContent")
	}

	var userIds []int64
	if req.All {
		allUserIds, err := fetchAllUserIds()
		if err != nil {
			log.Printf("Ошибка извлечения всех ID пользователей: %v", err)
			return nil, fmt.Errorf("ошибка сервера")
		}
		userIds = allUserIds
		log.Printf("Отправляем сообщение всем пользователям. Всего ID: %d", len(userIds))
	} else {
		userIds = req.UserIds
		log.Printf("Отправляем сообщение указанным пользователям. ID: %v", userIds)
	}

	successCount := 0
	failedUserIds := []int64{}

	for _, userID := range userIds {
		reqCopy := *req
		reqCopy.UserIds = []int64{userID}

		msgData, err := proto.Marshal(&reqCopy)
		if err != nil {
			log.Printf("Ошибка сериализации запроса для пользователя %d: %v", userID, err)
			failedUserIds = append(failedUserIds, userID)
			continue
		}

		pubAck, err := m.js.Publish("MESSAGES.send_message", msgData)
		if err != nil || pubAck == nil || pubAck.Sequence == 0 {
			log.Printf("Ошибка публикации для пользователя %d: %v", userID, err)
			failedUserIds = append(failedUserIds, userID)
			continue
		}

		successCount++
		log.Printf("Сообщение успешно отправлено пользователю %d", userID)
	}

	response := &proto_tg_service.SendMessageResponse{
		Success:       len(failedUserIds) == 0,
		SuccessCount:  int32(successCount),
		FailureCount:  int32(len(failedUserIds)),
		FailedUserIds: failedUserIds,
	}

	return response, nil
}

type UserID struct {
	UserID int64 `gorm:"column:user_id"`
}

func (UserID) TableName() string {
	return "users"
}

// fetchAllUserIds - вспомогательная функция для постраничного извлечения всех ID пользователей
func fetchAllUserIds() ([]int64, error) {
	db := database.GetDB()

	const batchSize = 1000
	var allUserIds []int64
	offset := 0

	for {
		var users []UserID
		err := db.Limit(batchSize).Offset(offset).Find(&users).Error
		if err != nil {
			return nil, fmt.Errorf("ошибка запроса к базе данных: %v", err)
		}

		if len(users) == 0 {
			break // завершение, если данных больше нет
		}

		for _, user := range users {
			allUserIds = append(allUserIds, user.UserID)
		}

		offset += batchSize
	}

	log.Printf("Всего пользователей для рассылки: %d", len(allUserIds))
	return allUserIds, nil
}
