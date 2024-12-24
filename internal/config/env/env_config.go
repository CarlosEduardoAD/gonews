package env

import (
	"log"
	"os"

	"github.com/CarlosEduardoAD/go-news/internal/config/logging"
	"github.com/joho/godotenv"
)

func GetEnv(key, fallback string) string {
	logger := logging.GenerateLogrus()

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Erro ao carregar o arquivo .env")
			logger.Printf("Erro ao carregar o arquivo .env: %v", err)
		}
	}

	// Buscar a variável de ambiente
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
