package storage

import (
	"context"
	"encoding/json"
	"strings"

	models "github.com/Konscig/foodelivery-pet/internal/storage/models"
)

func (s *PGstorage) UpdateOrderStatus(id string, status models.OrderStatus) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, id)
	return err
}

func (s *PGstorage) UpdateDeliveryStatus(id string, status string) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE deliveries SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, id)
	return err
}

func (s *PGstorage) UpdateRestaurantStats(restaurantID string) error {
	var reviews []models.Review
	rows, err := s.db.Query(context.Background(), `
		SELECT rating, comment FROM reviews WHERE restaurant_id = $1
	`, restaurantID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sum int32
	var comments []string
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.Rating, &review.Comment); err != nil {
			return err
		}
		sum += review.Rating
		comments = append(comments, review.Comment)
	}

	avg := float64(sum) / float64(len(reviews))
	wordCloud := BuildWordCloud(strings.Join(comments, " "))

	wordCloudJSON, err := json.Marshal(wordCloud)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(context.Background(), `
		INSERT INTO restaurant_stats (restaurant_id, avg_rating, reviews_count, word_cloud_json)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (restaurant_id) DO UPDATE
		SET avg_rating = $2, reviews_count = $3, word_cloud_json = $4
	`, restaurantID, avg, len(reviews), string(wordCloudJSON))
	return err
}

func BuildWordCloud(input string) map[string]int {
	// Placeholder implementation
	return map[string]int{"example": 1}
}
