package main

import (
	"log"

	"telegram-grpc-service/pkg/worker"
)

func main() {
	log.Println("Воркер запущен")
	if err := worker.StartWorker(); err != nil {
		log.Fatalf("Ошибка при запуске воркера: %v", err)
	}
}
