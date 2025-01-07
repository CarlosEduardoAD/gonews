package jobmodel

import (
	"fmt"
	"testing"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func TestInvalidVerifyMonday(t *testing.T) {
	sendEmailJob := NewSendEmailJob("1", "test@gmail.com", int64(time.Now().Second()), "")

	err := sendEmailJob.VerifyMonday()

	assert.NotNil(t, err)
}

func TestValidVerifyMonday(t *testing.T) {
	monday := utils.ReturnNextMonday()
	sendEmailJob := NewSendEmailJob("1", "test@gmail.com", monday, "")

	err := sendEmailJob.VerifyMonday()

	assert.Nil(t, err)
}

func TestInvalidAddAndEnqueueTask(t *testing.T) {
	sendEmailJob := NewSendEmailJob("", "", int64(time.Now().Second()), "")

	err := sendEmailJob.AddAndEnqueueTask(nil)

	assert.NotNil(t, err)
}

func TestValidAddAndEnqueueTask(t *testing.T) {

	password := env.GetEnv("REDIS_PASSWORD", "redis-password")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "gonews-redis:6379", Password: password})

	defer client.Close()

	sendEmailJob := NewSendEmailJob(utils.GenerateRandomString(8), fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)), int64(time.Now().Second()), "send_email")

	err := sendEmailJob.AddAndEnqueueTask(client)

	assert.Nil(t, err)
}
