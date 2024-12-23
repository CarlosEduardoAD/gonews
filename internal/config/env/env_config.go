package env

import (
	"os"

	"github.com/CarlosEduardoAD/go-news/internal/config/logging"
	"github.com/joho/godotenv"
)

func GetEnv(key, fallback string) string {
	logger := logging.GenerateLogrus()

	// Carregar o arquivo .env somente em ambiente local (opcional)
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			logger.Fatalf("Erro ao carregar o arquivo .env: %v", err)
		}
	}

	// Buscar a variável de ambiente
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
