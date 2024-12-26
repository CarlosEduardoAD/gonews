package jobmodel

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/models/consts"
	"github.com/hibiken/asynq"
)

type SendEmailJob struct {
	Id    string    `json:"id"`
	Email string    `json:"email"`
	TTD   time.Time `json:"ttd"` // Time To Deliver
	Type  string    `json:"type"`
}

func NewSendEmailJob(id, email string, ttd time.Time, task_type string) *SendEmailJob {
	return &SendEmailJob{
		Id:    id,
		Email: email,
		TTD:   ttd,
		Type:  task_type,
	}
}

func (sej *SendEmailJob) validate() error {
	if sej.Id == "" {
		return errors.New("invalid id")
	}

	if sej.Email == "" {
		return errors.New("invalid email")
	}

	if sej.TTD.IsZero() {
		return errors.New("invalid time to deliver")
	}

	if sej.Type == "" || sej.Type != consts.SendEmail {
		return errors.New("invalid type")
	}

	return nil

}

func (sej *SendEmailJob) VerifyMonday() error {
	if int(sej.TTD.Weekday()) != 1 {
		return errors.New("newsletter can only be delivered on mondays")
	}

	return nil
}

func (sej *SendEmailJob) AddAndEnqueueTask(task_manager *asynq.Client) error {
	err := sej.validate()

	if err != nil {
		return err
	}

	payload := map[string]interface{}{"id": sej.Id, "email": sej.Email, "ttd": sej.TTD}
	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	task := asynq.NewTask(sej.Type, payloadBytes)

	_, err = task_manager.Enqueue(task, asynq.ProcessIn(1*time.Minute), asynq.MaxRetry(2), asynq.Timeout(20*time.Second))

	if err != nil {
		return err
	}

	return nil
}
