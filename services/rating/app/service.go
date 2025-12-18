package app

import (
	"log"

	"github.com/Konscig/foodelivery-pet/services/rating/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RatingService struct {
	db        *gorm.DB
	publisher *Publisher
}

func NewRatingService(db *gorm.DB, publisher *Publisher) *RatingService {
	return &RatingService{
		db:        db,
		publisher: publisher,
	}
}

func (s *RatingService) SetOrderRated(orderID, userID string, rating uint8, comment string) error {
	r := models.Rating{
		ID:      uuid.NewString(),
		OrderID: orderID,
		UserID:  userID,
		Rating:  rating,
		Comment: comment,
	}

	if err := s.db.Create(&r).Error; err != nil {
		return err
	}

	if err := s.publisher.PublishOrderRated(orderID, rating, comment); err != nil {
		log.Println("Failed to publish order.rated:", err)
	}

	// TODO: обработка текста для облака слов
	// Например: подсчёт ключевых слов, сохранение в Redis

	return nil
}
