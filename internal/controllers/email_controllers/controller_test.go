package email_controllers

import (
	"fmt"
	"testing"

	"github.com/CarlosEduardoAD/go-news/internal/config/db"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCheckInEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})

	ec := NewEmailController(psql_db)

	email := emailmodel.NewEmailModel(fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)))

	link, err := ec.CheckInEmail(email)

	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.NotEmpty(t, link)
}

func TestFailToAuthorizeEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})
	valid_token := "invalid-token"

	ec := NewEmailController(psql_db)

	err := ec.AuthorizeEmail(valid_token)

	assert.NotNil(t, err)
}

func TestAuthorizeEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})
	valid_token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsc2tyby10ZXN0LWVtYWlsQHRlc3QuY29tIiwiZXhwIjoxNzM0NDU1MzY4fQ.dA_G9uq9zzT82FjAzzVXILSzyAmYtvsFMU2yoscpZY8"

	ec := NewEmailController(psql_db)

	err := ec.AuthorizeEmail(valid_token)

	assert.Nil(t, err)
}

func TestInvalidResendEmail(t *testing.T) {
	controller := NewEmailController(nil)
	email := "alskro-test-email@test.com"
	err := controller.ResendEmail(email)

	assert.NotNil(t, err)
}

func TestResendEmail(t *testing.T) {
	session := db.GenereateDB()
	session.AutoMigrate(&emailmodel.EmailModel{})
	controller := NewEmailController(session)
	token, err := shared.GenerateToken(jwt.MapClaims{
		"email": "alskro-test-email@test.com",
	})

	assert.Nil(t, err)

	err = controller.ResendEmail(token)

	assert.Nil(t, err)
}

func TestDismissEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})

	ec := NewEmailController(psql_db)
	email := "alskro-test-email@test.com"

	err := ec.DismissEmail(email)

	assert.Nil(t, err)
}
