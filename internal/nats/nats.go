package nats

import (
	"github.com/nats-io/nats.go"
)

func ConnectToNUTS() (*nats.Conn, error) {
	natsURL := "nats://172.18.0.2:4222"
	if natsURL == "" {
		natsURL = nats.DefaultURL // Используем URL по умолчанию, если переменная окружения пустая
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
