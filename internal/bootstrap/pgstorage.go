package bootstrap

import (
	"fmt"
	"log"

	"github.com/Konscig/foodelivery-pet/config"
	"github.com/Konscig/foodelivery-pet/internal/storage"
)

func InitPGStorage(cfg *config.Config) storage.Storage {

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	storage, err := storage.NewPGStorage(connectionString)
	if err != nil {
		log.Panicf("ошибка инициализации БД, %v", err)
	}
	return storage
}
