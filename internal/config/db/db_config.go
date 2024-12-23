package db

import (
	"fmt"
	"sync"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once     sync.Once
	instance *gorm.DB
)

// Return a singleton from the database connection
func GenereateDB() *gorm.DB {
	once.Do(func() {
		host := env.GetEnv("POSTGRES_HOST", "db")
		password := env.GetEnv("POSTGRES_USER", "admin")
		user := env.GetEnv("POSTGRES_PASSWORD", "Admin")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=gonews port=5432 sslmode=disable", host, user, password)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}

		instance = db
	})

	return instance
}
