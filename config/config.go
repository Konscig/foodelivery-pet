package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
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
}

func Load() (*Config, error) {
	godotenv.Load(".env")

	cfg := &Config{}

	cfg.Database.Host = os.Getenv("PG_HOST")
	cfg.Database.Port = os.Getenv("PG_PORT")
	cfg.Database.Username = os.Getenv("PG_USER")
	cfg.Database.Password = os.Getenv("PG_PASSWORD")
	cfg.Database.DBName = os.Getenv("PG_DB")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")

	cfg.Kafka.Broker = os.Getenv("KAFKA_BROKER")

	return cfg, nil
}
