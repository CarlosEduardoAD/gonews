package jobcontrollers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/CarlosEduardoAD/go-news/internal/config/db"
	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	newsapi_client "github.com/CarlosEduardoAD/go-news/internal/config/news-api"
	"github.com/CarlosEduardoAD/go-news/internal/config/task_queue"
	emailmodel "github.com/CarlosEduardoAD/go-news/internal/models/email_model"
	jobmodel "github.com/CarlosEduardoAD/go-news/internal/models/job_model"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/barthr/newsapi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hibiken/asynq"
)

type JobController struct {
	Client *asynq.Client
}

// required: client
func NewJobController(client *asynq.Client) *JobController {
	return &JobController{
		Client: client,
	}
}

func (jc *JobController) CreateTask(task *jobmodel.SendEmailJob) error {
	if jc.Client == nil {
		return errors.New("you've created a controller without a valid connection, this is only allowed if you want to execute a task")
	}

	err := task.AddAndEnqueueTask(jc.Client)

	if err != nil {
		return err
	}

	return nil
}

func (jc *JobController) ExecuteTask(ctx context.Context, t *asynq.Task) error {
	var emailJobPayload jobmodel.SendEmailJob
	var email_model emailmodel.EmailModel

	if err := json.Unmarshal(t.Payload(), &emailJobPayload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	session := db.GenereateDB()

	newsApiClient := newsapi_client.GenerateNewsApi()

	articleResponse, err := newsApiClient.GetEverything(ctx, &newsapi.EverythingParameters{Keywords: "golang"})

	if err != nil {
		return err
	}

	email, err := email_model.SelectOneByEmail(session, emailJobPayload.Email)

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

	unsubscribeToken, err := shared.GenerateToken(jwt.MapClaims{"email": emailJobPayload.Email})

	if err != nil {
		return err
	}

	unsubscribeLink := fmt.Sprintf("http://localhost:3000/api/v1/emails/dismiss?token=%s", unsubscribeToken)

	template, err := utils.LoadTemplate("internal/views/templates/news.html", utils.EmailData{NewsTitle: articleOfTheDay.Title, NewsDescription: articleOfTheDay.Description, NewsLink: articleOfTheDay.URL, UnsubscribeLink: unsubscribeLink})

	if err != nil {
		return err
	}

	err = email_sender.SendEmail(emailJobPayload.Email, "Sua newsletter de go de hoje chegou!", template)

	if err != nil {
		return err
	}

	client := task_queue.GenerateAsynqClient()

	defer client.Close()

	email_job := jobmodel.NewSendEmailJob(emailJobPayload.Id, emailJobPayload.Email, utils.ReturnNextMonday(), "send_email")
	err = email_job.AddAndEnqueueTask(client)

	if err != nil {
		return err
	}
	return nil
}
