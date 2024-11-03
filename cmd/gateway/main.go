package gateway

import (
	"context"
	"log"
	"net/http"
	"os"

	proto_tg_service "github.com/0sokrat0/telegram-grpc-service/gen/go/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func MainGateway() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcServerEndpoint := os.Getenv("GRPC_SERVER_ENDPOINT")
	if grpcServerEndpoint == "" {
		grpcServerEndpoint = "localhost:50051"
	}

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := proto_tg_service.RegisterMessagingServiceHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Ошибка при регистрации сервиса: %v", err)
	}

	log.Println("Запуск HTTP сервера на порту 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка при запуске HTTP сервера: %v", err)
	}
}
