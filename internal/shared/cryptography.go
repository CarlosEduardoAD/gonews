package shared

import (
	"errors"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/golang-jwt/jwt/v5"
)

// TODO: Realocar essa struct, ela não pertence aqui
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(claims jwt.MapClaims) (string, error) {
	secretKey := []byte(env.GetEnv("SECRET_KEY", "my-secret-key"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", GenerateError(err)
	}

	return tokenString, nil
}

func CompareTokenAndReturnClaims(stringToken string) (jwt.Claims, error) {
	secretKey := []byte(env.GetEnv("SECRET_KEY", "my-secret-key"))
	result, err := jwt.ParseWithClaims(stringToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, GenerateError(errors.New("invalid-token"))
	}

	return result.Claims, nil

}
