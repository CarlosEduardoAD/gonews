package db

import (
	"sync"

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
		dsn := "host=db user=admin password=admin dbname=gonews port=5432 sslmode=disable"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}

		instance = db
	})

	return instance
}
