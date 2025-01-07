package jobcontrollers

import (
	"fmt"
	"testing"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	"github.com/CarlosEduardoAD/go-news/internal/models/consts"
	jobmodel "github.com/CarlosEduardoAD/go-news/internal/models/job_model"
	"github.com/CarlosEduardoAD/go-news/internal/utils"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func TestInvalidCreateTask(t *testing.T) {
	password := env.GetEnv("REDIS_PASSWORD", "redis-password")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "gonews-redis:6379", Password: password})

	controller := NewJobController(client)
	myTask := jobmodel.NewSendEmailJob("", "", int64(time.Now().Second()), "")

	err := controller.CreateTask(myTask)

	assert.NotNil(t, err)
}

func TestValidCreateTask(t *testing.T) {
	password := env.GetEnv("REDIS_PASSWORD", "redis-password")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "gonews-redis:6379", Password: password})

	controller := NewJobController(client)
	email := fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8))
	myTask := jobmodel.NewSendEmailJob(utils.GenerateRandomString(8), email, utils.ReturnNextMonday(), consts.SendEmail)

	err := controller.CreateTask(myTask)

	assert.Nil(t, err)
}

// func TestExecuteTask(t *testing.T) {
// 	controller := NewJobController(nil)
// 	err := controller.ExecuteTask()

// 	assert.NotNil(t, err)
// }

// Esse vai ser  mais díficil de testar, pois ele depende de uma conexão com o servidor do asynq
