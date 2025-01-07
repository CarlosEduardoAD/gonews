package jobcontrollers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/CarlosEduardoAD/go-news/internal/config/db"
	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	newsapi_client "github.com/CarlosEduardoAD/go-news/internal/config/news-api"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	jobmodel "github.com/CarlosEduardoAD/go-news/internal/models/job_model"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/barthr/newsapi"
	"github.com/gocraft/work"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JobController struct{}

// required: client
func NewJobController() *JobController {
	return &JobController{}
}

func (jc *JobController) CreateTask(task *jobmodel.SendEmailJob) error {

	err := task.AddAndEnqueueTask()

	if err != nil {
		return err
	}

	return nil
}

func (jc *JobController) ExecuteTask(job *work.Job) error {
	var email_model emailmodel.EmailModel

	payloadEmail := job.ArgString("email")

	session := db.GenereateDB()

	newsApiClient := newsapi_client.GenerateNewsApi()

	articleResponse, err := newsApiClient.GetEverything(context.Background(), &newsapi.EverythingParameters{Keywords: "golang"})

	if err != nil {
		return err
	}

	email, err := email_model.SelectOneByEmail(session, payloadEmail)

	if err != nil {
		return err
	}

	if !email.Authorized {

		return errors.New("unauthorized")
	}

	articleOfTheDay := articleResponse.Articles[0]

	host := env.GetEnv("MAILTRAP_HOST", "my-host")
	port, err := strconv.Atoi(env.GetEnv("MAILTRAP_PORT", "my-port"))

	if err != nil {
		panic(err)
	}

	username := env.GetEnv("MAILTRAP_USERNAME", "my-user")
	password := env.GetEnv("MAILTRAP_PASSWORD", "my-password")

	email_sender := shared.GenerateEmailSender(host, port, username, password)

	unsubscribeToken, err := shared.GenerateToken(jwt.MapClaims{"email": payloadEmail})

	if err != nil {
		return err
	}

	unsubscribeLink := fmt.Sprintf("http://localhost:3000/api/v1/emails/dismiss?token=%s", unsubscribeToken)

	template, err := utils.LoadTemplate("internal/views/templates/news.html", utils.EmailData{NewsTitle: articleOfTheDay.Title, NewsDescription: articleOfTheDay.Description, NewsLink: articleOfTheDay.URL, UnsubscribeLink: unsubscribeLink})

	if err != nil {
		return err
	}

	err = email_sender.SendEmail(payloadEmail, "Sua newsletter de go de hoje chegou!", template)

	if err != nil {
		return err
	}

	email_job := jobmodel.NewSendEmailJob(uuid.NewString(), payloadEmail, utils.ReturnNextMonday(), "send_email")
	err = email_job.AddAndEnqueueTask()

	if err != nil {
		return err
	}
	return nil
}
