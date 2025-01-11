package email_controllers

import (
	"fmt"
	"log"
	"testing"

	"github.com/CarlosEduardoAD/go-news/internal/config/db"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Mock da função LoadTemplate para testes

func TestCheckInEmail(t *testing.T) {
	psql_db := db.GenereateDBTest()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})

	ec := NewEmailController(psql_db)

	email := emailmodel.NewEmailModel(fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)))

	link, err := ec.CheckInEmail(email.Email)

	log.Println(err)

	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.NotEmpty(t, link)
}

func TestFailToAuthorizeEmail(t *testing.T) {
	psql_db := db.GenereateDBTest()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})
	valid_token := "invalid-token"

	ec := NewEmailController(psql_db)

	err := ec.AuthorizeEmail(valid_token)

	assert.NotNil(t, err)
}

func TestAuthorizeEmail(t *testing.T) {
	psql_db := db.GenereateDBTest()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})

	ec := NewEmailController(psql_db)
	email := emailmodel.NewEmailModel(fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)))

	link, err := ec.CheckInEmail(email.Email)

	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.NotEmpty(t, link)

	err = ec.AuthorizeEmail(link)

	assert.Nil(t, err)
}

func TestInvalidResendEmail(t *testing.T) {
	controller := NewEmailController(nil)
	email := "alskro-test-email@test.com"
	err := controller.ResendEmail(email)

	assert.NotNil(t, err)
}

func TestResendEmail(t *testing.T) {
	psql_db := db.GenereateDBTest()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})
	controller := NewEmailController(psql_db)
	email := emailmodel.NewEmailModel(fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)))

	link, err := controller.CheckInEmail(email.Email)

	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.NotEmpty(t, link)

	token, err := shared.GenerateToken(jwt.MapClaims{
		"email": email.Email,
	})

	assert.Nil(t, err)

	err = controller.ResendEmail(token)

	assert.Nil(t, err)
}

func TestDismissEmail(t *testing.T) {
	psql_db := db.GenereateDBTest()
	psql_db.AutoMigrate(&emailmodel.EmailModel{})

	ec := NewEmailController(psql_db)

	email := emailmodel.NewEmailModel(fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)))
	_, err := ec.CheckInEmail(email.Email)

	assert.Nil(t, err)

	err = ec.DismissEmail(email.Email)

	assert.Nil(t, err)
}
