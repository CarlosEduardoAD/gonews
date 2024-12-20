package shared

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func GenerateValidError(t *testing.T) {
	err := GenerateError(errors.New("ERROR MESSAGE EXAMPLE"))
	assert.NotNil(t, err)
}

func TestJWTCreationWithValidClaims(t *testing.T) {
	validClaims := jwt.MapClaims{
		"email": fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := GenerateToken(validClaims)

	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token)
}

func TestInvalidJWTComparison(t *testing.T) {
	// There's no problem in using exposing this token
	// since it doesn't have any sensitive information
	valid_token := "invalid-token"
	token, err := CompareTokenAndReturnClaims(valid_token)

	assert.NotNil(t, err)
	assert.Nil(t, token)
}

func TestValidJWTComparison(t *testing.T) {
	// There's no problem in using exposing this token
	// since it doesn't have any sensitive information
	valid_token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImthcmwuZGV2Y29udGF0b0BnbWFpbC5jb20iLCJleHAiOjE3MzQ0NTQxNTZ9.u7I1aKOo37rral4x1ZdpzlqA8xeKBhYP3zKdT1XmXGw"
	token, err := CompareTokenAndReturnClaims(valid_token)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestInvalidSendEmail(t *testing.T) {
	host := env.GetEnv("MAILTRAP_HOST", "my-host")
	port, err := strconv.Atoi(env.GetEnv("MAILTRAP_PORT", "my-port"))

	if err != nil {
		panic(err)
	}

	username := env.GetEnv("MAILTRAP_USERNAME", "my-user")
	password := env.GetEnv("MAILTRAP_PASSWORD", "my-password")

	email_sender := GenerateEmailSender(host, port, username, password)
	err = email_sender.SendEmail("", "", "")

	assert.NotNil(t, err)
}

func TestValidSendEmail(t *testing.T) {
	host := env.GetEnv("MAILTRAP_HOST", "my-host")
	port, err := strconv.Atoi(env.GetEnv("MAILTRAP_PORT", "my-port"))

	if err != nil {
		panic(err)
	}

	username := env.GetEnv("MAILTRAP_USERNAME", "my-user")
	password := env.GetEnv("MAILTRAP_PASSWORD", "my-password")

	email_sender := GenerateEmailSender(host, port, username, password)
	err = email_sender.SendEmail("karl.devcontato@gmail.com", "Sua newsletter de go de hoje chegou!", "<h1>Notícia de Hoje: Go vai recebe atualização!</h1>")

	assert.Nil(t, err)
}
