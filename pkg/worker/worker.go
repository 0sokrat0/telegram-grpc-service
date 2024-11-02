package worker

import (
	"fmt"
	"github.com/0sokrat0/telegram-grpc-service/config"
	"github.com/0sokrat0/telegram-grpc-service/pkg/telegram"
	"sync"
	"time"

	proto_tg_service "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

// StartWorker запускает воркер для обработки сообщений
func StartWorker() error {
	cfg := config.GetConfig()
	natsURL := fmt.Sprintf("%s:%d", cfg.NATSHost, cfg.NATSPort)

	// Подключение к NATS серверу
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к NATS: %v", err)
	}
	defer nc.Close()

	// Инициализация JetStream контекста
	js, err := nc.JetStream()
	if err != nil {
		return fmt.Errorf("не удалось инициализировать JetStream: %v", err)
	}

	// Создаем поток "MESSAGES", если он не существует
	_, err = js.StreamInfo("MESSAGES")
	if err != nil {
		if err == nats.ErrStreamNotFound {
			_, err = js.AddStream(&nats.StreamConfig{
				Name:     "MESSAGES",
				Subjects: []string{"MESSAGES.*"},
				Storage:  nats.FileStorage,
			})
			if err != nil {
				return fmt.Errorf("не удалось создать поток MESSAGES: %v", err)
			}
		} else {
			return fmt.Errorf("ошибка при получении информации о потоке MESSAGES: %v", err)
		}
	}

	// Подписываемся на поток "MESSAGES" с темой "MESSAGES.send_message"
	sub, err := js.PullSubscribe("MESSAGES.send_message", "worker", nats.BindStream("MESSAGES"))
	if err != nil {
		return fmt.Errorf("не удалось подписаться на поток MESSAGES: %v", err)
	}

	// Семафор для ограничения одновременных отправок (до 30 сообщений)
	sem := make(chan struct{}, 30) // 30 одновременных отправок

	for {
		// Извлекаем до 10 сообщений из очереди
		msgs, err := sub.Fetch(10)
		if err != nil {
			// Обработка ошибок извлечения
			time.Sleep(time.Second)
			continue
		}

		for _, msg := range msgs {
			// Десериализуем сообщение
			var req proto_tg_service.SendMessageRequest
			err = proto.Unmarshal(msg.Data, &req)
			if err != nil {
				// Логирование ошибки десериализации
				fmt.Printf("Ошибка десериализации сообщения: %v\n", err)
				msg.Nak() // Отмена подтверждения
				continue
			}

			var wg sync.WaitGroup

			for _, userID := range req.UserIds {
				wg.Add(1)
				go func(userID int64) {
					defer wg.Done()
					sem <- struct{}{} // Захватываем слот в семафоре

					// Отправляем сообщение пользователю
					err := telegram.SendMessageToUser(userID, req)

					if err != nil {
						// Логирование ошибки отправки
						fmt.Printf("Ошибка отправки сообщения пользователю %d: %v\n", userID, err)
					}

					<-sem // Освобождаем слот в семафоре
				}(userID)
			}

			// Ждем завершения отправки сообщений всем пользователям в запросе
			wg.Wait()

			// Подтверждаем успешную обработку сообщения
			msg.Ack()
		}
	}
}
