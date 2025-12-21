package main

import (
	"log"

	"github.com/Konscig/foodelivery-pet/internal/services/restaurant"
)

func main() {
	service := restaurant.NewService()
	if err := service.Start(); err != nil {
		log.Fatalf("Ошибка запуска сервиса: %v", err)
	}
}
