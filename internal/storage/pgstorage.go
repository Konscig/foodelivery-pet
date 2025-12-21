package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGstorage struct {
	db *pgxpool.Pool
}

func NewPGStorage(connString string) (Storage, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}
	storage := &PGstorage{
		db: db,
	}
	err = storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGstorage) initTables() error {
	_, err := s.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			rest_id VARCHAR(255) NOT NULL,
			status VARCHAR(50) NOT NULL,
			items JSONB,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return errors.Wrap(err, "init orders table")
	}

	_, err = s.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS deliveries (
			id VARCHAR(255) PRIMARY KEY,
			order_id VARCHAR(255) NOT NULL,
			courier_id VARCHAR(255),
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return errors.Wrap(err, "init deliveries table")
	}

	_, err = s.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS reviews (
			id VARCHAR(255) PRIMARY KEY,
			order_id VARCHAR(255) NOT NULL,
			restaurant_id VARCHAR(255) NOT NULL,
			rating INT NOT NULL,
			comment TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return errors.Wrap(err, "init reviews table")
	}

	_, err = s.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS restaurant_stats (
			restaurant_id VARCHAR(255) PRIMARY KEY,
			avg_rating FLOAT,
			reviews_count INT DEFAULT 0,
			word_cloud_json TEXT
		)
	`)
	if err != nil {
		return errors.Wrap(err, "init restaurant_stats table")
	}

	return nil
}
