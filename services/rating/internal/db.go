package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	godotenv.Load(".env")

	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB_RATE")
	pgPort := os.Getenv("PG_PORT")
	pgHost := os.Getenv("PG_HOST")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgDB,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id UUID PRIMARY KEY,
            user_id UUID NOT NULL,
            restaurant_id UUID NOT NULL,
            status TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT now()
        ) PARTITION BY HASH (id);
    `)

	for i := 0; i < 4; i++ {
		db.Exec(fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS orders_shard%d PARTITION OF orders
            FOR VALUES WITH (MODULUS 4, REMAINDER %d);
        `, i+1, i))
	}

	return db, nil
}
