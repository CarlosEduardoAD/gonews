package db

import (
	"fmt"
	"log"
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
		password := env.GetEnv("POSTGRES_PASSWORD", "admin")
		user := env.GetEnv("POSTGRES_USER", "admin")
		database := env.GetEnv("POSTGRES_DB", "gonews")

		sslmode := "disable"

		if database == "railway" {
			sslmode = "require"
		}

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s", host, user, password, database, sslmode)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Println("Erro ao conectar ao banco de dados")
			log.Fatal(err)
		}

		instance = db
	})

	return instance
}

func GenereateDBTest() *gorm.DB {
	once.Do(func() {
		host := env.GetEnv("POSTGRES_HOST", "db")
		password := env.GetEnv("POSTGRES_PASSWORD", "admin")
		user := env.GetEnv("POSTGRES_USER", "admin")
		database := "gonews_test"

		sslmode := "disable"

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s", host, user, password, database, sslmode)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Println("Erro ao conectar ao banco de dados")
			log.Fatal(err)
		}

		instance = db
	})

	return instance
}
