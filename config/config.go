package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN" env-required:"true"`
	GrpcHost string `env:"GRPC_HOST" env-default:"localhost"`
	GrpcPort int    `env:"GRPC_PORT" env-default:"50051"`
	NATSHost string `env:"NATS_HOST" env-default:"localhost"`
	NATSPort int    `env:"NATS_PORT" env-default:"50052"`
}

var instance *Config
var once sync.Once

// GetConfig возвращает экземпляр конфигурации, загруженный один раз
func GetConfig() *Config {
	once.Do(func() {
		// Загрузка из файла .env, если он существует
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Файл .env не найден, загрузка из переменных окружения")
		}

		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			log.Fatalf("Ошибка загрузки конфигурации: %v", err)
		}
	})
	return instance
}
