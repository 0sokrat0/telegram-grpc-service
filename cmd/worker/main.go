package worker

import (
	"github.com/0sokrat0/telegram-grpc-service/internal/worker"
	"log"
)

func MainWorker() {
	log.Println("Воркер запущен")
	if err := worker.StartWorker(); err != nil {
		log.Fatalf("Ошибка при запуске воркера: %v", err)
	}
}
