package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/0sokrat0/telegram-grpc-service/cmd/gateway"
	"github.com/0sokrat0/telegram-grpc-service/cmd/server"
	"github.com/0sokrat0/telegram-grpc-service/cmd/worker"
)

func main() {
	var wg sync.WaitGroup

	// Запуск gRPC сервера
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Запуск gRPC сервера...")
		server.MainServer()
	}()

	// Запуск воркера
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Запуск воркера...")
		worker.MainWorker()
	}()

	// Запуск HTTP сервера (gateway)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Запуск HTTP сервера (gateway)...")
		gateway.MainGateway()
	}()

	// Обработка сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("Завершается работа всех сервисов...")
	wg.Wait()
}
