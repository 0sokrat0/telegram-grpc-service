package nats

import (
	"github.com/nats-io/nats.go"
	"os"
)

func ConnectToNUTS() (*nats.Conn, error) {
	natsURL := os.Getenv("NATS_URL")
	natsPort := os.Getenv("NATS_PORT")
	natsFull := "nats://" + natsURL + ":" + natsPort

	nc, err := nats.Connect(natsFull)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
