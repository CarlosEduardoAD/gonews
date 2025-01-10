package emailmodel

import (
	"errors"
	"net/mail"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/models/consts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type EmailModel struct {
	Id         string    `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"uniqueIndex:unique_email,sort:desc"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Authorized bool      `json:"authorized"`
	Deleted    bool      `json:"deleted"`
}

type SaveEmailModelDTO struct {
	Email string `json:"email"`
}

func NewEmailModel(email string) *EmailModel {
	return &EmailModel{
		Id:         uuid.New().String(),
		Email:      email,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Deleted:    false,
		Authorized: false,
	}

}

func (em *EmailModel) CheckInvalidEmail() error {
	if _, err := mail.ParseAddress(em.Email); err != nil {
		return err
	}

	return nil
}

func (em *EmailModel) Create(session *gorm.DB) error {
	if em.CheckInvalidEmail() != nil {
		return errors.New(consts.InvalidEmail)
	}

	ctx := session.Create(em)
	if ctx.Error != nil {
		if pgError, ok := ctx.Error.(*pgconn.PgError); ok {
			if pgError.Code == "23505" {
				return errors.New("email already exists")
			}
		}
		return ctx.Error
	}

	return nil
}

func (em *EmailModel) SelectOne(session *gorm.DB, id string) (*EmailModel, error) {
	var email EmailModel

	result := session.Where("id = ?", id).First(&email)
	if result.Error != nil {

		return nil, result.Error
	}

	return &email, nil
}

func (em *EmailModel) SelectOneByEmail(session *gorm.DB, search_email string) (*EmailModel, error) {
	var email EmailModel

	result := session.Where("email = ?", search_email).First(&email)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	return &email, nil
}

func (em *EmailModel) Update(session *gorm.DB) error {
	result := session.Model(em).Updates(em)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// TODO: Paliativo
func (em *EmailModel) DismissEmail(session *gorm.DB) error {
	result := session.Model(em).Update("authorized", em.Authorized)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (em *EmailModel) DeleteEmail(session *gorm.DB, id string) error {
	ctx := session.Model(&EmailModel{}).Where("id = ?", id).Update("deleted", true)
	if ctx.Error != nil {
		return ctx.Error
	}

	return nil
}
