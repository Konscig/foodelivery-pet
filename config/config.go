package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		DBName   string
	}
	Redis struct {
		Host string
		Port string
	}
	Kafka struct {
		Broker string
	}
	GRPC struct {
		DeliveryPort   int
		OrderPort      int
		RatingPort     int
		RestaurantPort int
	}
}

func Load() (*Config, error) {
	godotenv.Load(".env")

	cfg := &Config{}

	cfg.Database.Host = os.Getenv("PG_HOST")
	cfg.Database.Port = getEnvAsInt("PG_PORT", 5432)
	cfg.Database.Username = os.Getenv("PG_USER")
	cfg.Database.Password = os.Getenv("PG_PASSWORD")
	cfg.Database.DBName = os.Getenv("PG_DB")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")

	cfg.Kafka.Broker = os.Getenv("KAFKA_BROKER")

	// gRPC ports (по умолчанию 50051, 50052, 50053, 50054)
	cfg.GRPC.DeliveryPort = getEnvAsInt("GRPC_DELIVERY_PORT", 50051)
	cfg.GRPC.OrderPort = getEnvAsInt("GRPC_ORDER_PORT", 8081)
	cfg.GRPC.RatingPort = getEnvAsInt("GRPC_RATING_PORT", 50053)
	cfg.GRPC.RestaurantPort = getEnvAsInt("GRPC_RESTAURANT_PORT", 50054)

	return cfg, nil
}

func getEnvAsInt(name string, defaultVal int) int {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	var i int
	_, err := fmt.Sscanf(val, "%d", &i)
	if err != nil {
		return defaultVal
	}
	return i
}
