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
	rows, err := s.db.Query(context.Background(), `
		SELECT rating, comment FROM reviews WHERE restaurant_id = $1
	`, restaurantID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sum int32
	var comments []string
	count := 0
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.Rating, &review.Comment); err != nil {
			return err
		}
		sum += review.Rating
		comments = append(comments, review.Comment)
		count++
	}

	var avg float64
	if count > 0 {
		avg = float64(sum) / float64(count)
	}

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
	`, restaurantID, avg, count, string(wordCloudJSON))
	return err
}

func (s *PGstorage) GetRestaurantStats(restaurantID string) (*models.RestaurantStats, error) {
	row := s.db.QueryRow(context.Background(), `
		SELECT avg_rating, reviews_count, word_cloud_json
		FROM restaurant_stats
		WHERE restaurant_id = $1
	`, restaurantID)
	var avg float64
	var count int
	var wordCloudJSON string
	err := row.Scan(&avg, &count, &wordCloudJSON)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return &models.RestaurantStats{RestaurantID: restaurantID, AvgRating: 0, ReviewsCount: 0, WordCloudJSON: "{}"}, nil
		}
		return nil, err
	}
	return &models.RestaurantStats{RestaurantID: restaurantID, AvgRating: avg, ReviewsCount: uint32(count), WordCloudJSON: wordCloudJSON}, nil
}

func BuildWordCloud(input string) map[string]int {
	// Placeholder implementation
	return map[string]int{"example": 1}
}
