package jobmodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/CarlosEduardoAD/go-news/internal/models/consts"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/hibiken/asynq"
)

var host = env.GetEnv("REDIS_HOST", "gonews-redis")
var password = env.GetEnv("REDIS_PASSWORD", "Carloseduardo08#")

var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp",
			fmt.Sprintf("%s:6379", host),
			redis.DialPassword(password))
	},
}

// Make an enqueuer with a particular namespace
var enqueuer = work.NewEnqueuer("go_news_namespace", redisPool)

type SendEmailJob struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	TTD   int64  `json:"ttd"` // Time To Deliver
	Type  string `json:"type"`
}

func NewSendEmailJob(id, email string, ttd int64, task_type string) *SendEmailJob {
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

	if sej.TTD == 0 {
		return errors.New("invalid time to deliver")
	}

	if sej.Type == "" || sej.Type != consts.SendEmail {
		return errors.New("invalid type")
	}

	return nil

}

func (sej *SendEmailJob) VerifyMonday() error {
	if sej.TTD != utils.ReturnNextMonday() {
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

	log.Println("payloadBytes: ", payloadBytes)

	_, err = enqueuer.EnqueueIn("send_email", int64(time.Second*15), work.Q{"id": sej.Id, "email": sej.Email, "ttd": sej.TTD})

	if err != nil {
		return err
	}

	return nil
}
