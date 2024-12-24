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

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=railway port=5432 sslmode=require", host, user, password)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Println("Erro ao conectar ao banco de dados")
			log.Fatal(err)
		}

		instance = db
	})

	return instance
}
