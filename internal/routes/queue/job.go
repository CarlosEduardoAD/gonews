package queue

import (
	"context"

	"github.com/CarlosEduardoAD/go-news/internal/config/env"
	jobcontrollers "github.com/CarlosEduardoAD/go-news/internal/controllers/job_controllers"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

func HandleEmailDelivery(ctx context.Context, t *asynq.Task) error {
	password := env.GetEnv("REDIS_PASSWORD", "redis-password")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "gonews-redis:6379", Password: password})

	controller := jobcontrollers.NewJobController(client)
	err := controller.ExecuteTask(ctx, t)

	if err != nil {
		logrus.Error(shared.GenerateError(err))
		return err
	}

	return nil
}
