package emailmodel

import (
	"fmt"
	"testing"

	"github.com/CarlosEduardoAD/go-news/internal/config/db"
	"github.com/CarlosEduardoAD/go-news/internal/models/consts"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvalidEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&EmailModel{})

	email_model := NewEmailModel("invalid_email")
	err := email_model.Create(psql_db)

	assert.NotNil(t, err)
	assert.EqualError(t, err, consts.InvalidEmail)
}

func TestCreateEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&EmailModel{})

	email_model := NewEmailModel(fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)))
	err := email_model.Create(psql_db)

	assert.Nil(t, err)
}

func TestSelectOneEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&EmailModel{})

	email_model := EmailModel{}
	result, err := email_model.SelectOne(psql_db, "1010d16b-3f43-486c-ae41-e7eb6b9f7717")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestSelectOneEmailByEmailString(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&EmailModel{})

	email_model := EmailModel{}
	result, err := email_model.SelectOneByEmail(psql_db, "alskro.devcontato@gmail.com")

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestAuthorizeEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&EmailModel{})

	email_model := EmailModel{}

	fetch, err := email_model.SelectOne(psql_db, "1010d16b-3f43-486c-ae41-e7eb6b9f7717")

	assert.NotNil(t, fetch)
	assert.Nil(t, err)

	fetch.Authorized = true
	err = fetch.Update(psql_db)

	assert.Nil(t, err)
}

func TestDismissEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&EmailModel{})

	email_model := EmailModel{}

	fetch, err := email_model.SelectOne(psql_db, "1010d16b-3f43-486c-ae41-e7eb6b9f7717")

	assert.NotNil(t, fetch)
	assert.Nil(t, err)

	fetch.Authorized = false

	err = fetch.DismissEmail(psql_db)

	assert.Nil(t, err)
}

func TestDeleteEmail(t *testing.T) {
	psql_db := db.GenereateDB()
	psql_db.AutoMigrate(&EmailModel{})

	email_model := EmailModel{}
	fetch, err := email_model.SelectOne(psql_db, "1010d16b-3f43-486c-ae41-e7eb6b9f7717")

	assert.NotNil(t, fetch)
	assert.Nil(t, err)

	err = fetch.DeleteEmail(psql_db, fetch.Id)

	assert.Nil(t, err)
}
