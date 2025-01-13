package email_controllers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	jobcontrollers "github.com/CarlosEduardoAD/go-news/internal/controllers/job_controllers"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	jobmodel "github.com/CarlosEduardoAD/go-news/internal/models/job_model"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmailController struct {
	db *gorm.DB
}

const (
	EmailNotFound = "email not found"
)

var (
	MAILTRAP_HOST     = env.GetEnv("MAILTRAP_HOST", "smtp.mailtrap.io")
	MAILTRAP_PORT     = env.GetEnv("MAILTRAP_PORT", "2525")
	MAILTRAP_USERNAME = env.GetEnv("MAILTRAP_USERNAME", "my-user")
	MAILTRAP_PASSWORD = env.GetEnv("MAILTRAP_PASSWORD", "my-password")
)

func NewEmailController(db *gorm.DB) *EmailController {
	return &EmailController{
		db,
	}
}

func (ec *EmailController) CheckInEmail(e string) (string, error) {

	server_url := env.GetEnv("SERVER_URL", "http://localhost:3000")

	token, err := shared.GenerateToken(jwt.MapClaims{
		"email": e,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	if err != nil {
		return "", shared.GenerateError(err)
	}

	doubleOptInLink := fmt.Sprintf("%s/emails/authorize?token=%s", server_url, token)

	port, err := strconv.Atoi(MAILTRAP_PORT)

	if err != nil {
		return "", err
	}

	email_sender := shared.GenerateEmailSender(MAILTRAP_HOST, port, MAILTRAP_USERNAME, MAILTRAP_PASSWORD)
	template, err := utils.LoadTemplate("internal/views/templates/check_in.html", utils.EmailData{ConfirmLink: doubleOptInLink})
	if err != nil {
		return "", err
	}

	err = email_sender.SendEmail(e, "Confirme seu email", template)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (ec *EmailController) AuthorizeEmail(token string) error {
	result, err := shared.CompareTokenAndReturnClaims(token)

	if err != nil {
		return err
	}

	claims := result.(*shared.Claims)

	email_model := emailmodel.NewEmailModel(claims.Email)
	err = email_model.Create(ec.db)

	if err != nil {
		return err
	}

	if email_model.Email == "" {
		return errors.New(EmailNotFound)
	}

	email_model.Authorized = true

	err = email_model.Update(ec.db)

	if err != nil {
		return err
	}

	job_controller := jobcontrollers.NewJobController()
	err = job_controller.CreateTask(jobmodel.NewSendEmailJob(email_model.Id, email_model.Email, utils.ReturnNextMonday(), "send_email"))

	if err != nil {
		log.Println(err)
		logrus.Error(shared.GenerateError(err))
		return err
	}

	port, err := strconv.Atoi(MAILTRAP_PORT)

	if err != nil {
		return err
	}

	email_sender := shared.GenerateEmailSender(MAILTRAP_HOST, port, MAILTRAP_USERNAME, MAILTRAP_PASSWORD)
	template, err := utils.LoadTemplate("internal/views/templates/confirmation.html", utils.EmailData{})
	if err != nil {
		return nil
	}

	err = email_sender.SendEmail(email_model.Email, "Email confirmado", template)

	if err != nil {
		return err
	}

	return nil
}

func (ec *EmailController) VerifyEmail(token string) error {
	email_model := emailmodel.EmailModel{}

	result, err := shared.CompareTokenAndReturnClaims(token)

	if err != nil {
		return err
	}

	claims := result.(*shared.Claims)

	fetch, err := email_model.SelectOneByEmail(ec.db, claims.Email)

	if err != nil {
		return err
	}

	if !fetch.Authorized {
		return errors.New(EmailNotFound)
	}

	return nil
}

func (ec *EmailController) ResendEmail(token string) error {
	email_model := emailmodel.EmailModel{}

	tokenClaims, err := shared.CompareTokenAndReturnClaims(token)

	if err != nil {
		return err
	}

	userClaims := tokenClaims.(*shared.Claims)

	fetch, err := email_model.SelectOneByEmail(ec.db, userClaims.Email)

	if err != nil {
		return err
	}

	if fetch == nil {
		return errors.New(EmailNotFound)
	}

	port, err := strconv.Atoi(MAILTRAP_PORT)

	if err != nil {
		return err
	}

	email_sender := shared.GenerateEmailSender(MAILTRAP_HOST, port, MAILTRAP_USERNAME, MAILTRAP_PASSWORD)

	new_token, err := shared.GenerateToken(jwt.MapClaims{
		"email": fetch.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	if err != nil {
		return err
	}

	confirmationLink := fmt.Sprintf("http://localhost:3000/api/v1/emails/authorize?token=%s", new_token)

	template, err := utils.LoadTemplate("internal/views/templates/check_in.html", utils.EmailData{
		ConfirmLink: confirmationLink,
	})

	if err != nil {
		return err
	}

	err = email_sender.SendEmail(fetch.Email, "Confirme seu email", template)

	if err != nil {
		return err
	}

	return nil
}

func (ec *EmailController) DismissEmail(token string) error {
	email_model := emailmodel.EmailModel{}

	tokenClaims, err := shared.CompareTokenAndReturnClaims(token)

	if err != nil {
		return err
	}

	userClaims := tokenClaims.(*shared.Claims)

	fetch, err := email_model.SelectOneByEmail(ec.db, userClaims.Email)

	if err != nil {
		return err
	}

	fetch.Authorized = false
	// FIXME: This wasn't suposed to be done, but Gorm just go nuts
	// when I try to update the "authorized" field using the common update function
	err = fetch.DismissEmail(ec.db)

	if err != nil {
		return err
	}

	return nil
}
