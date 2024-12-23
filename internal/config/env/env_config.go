package env

import (
	"log"
	"os"

	"github.com/CarlosEduardoAD/go-news/internal/config/logging"
	"github.com/joho/godotenv"
)

func GetEnv(key, fallback string) string {
	err := godotenv.Load("/app/.env")
	logger := logging.GenerateLogrus()

	if err != nil {
		log.Println("Erro ao carregar o arquivo .env")
		logger.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	value := os.Getenv(key)

	return value
}
