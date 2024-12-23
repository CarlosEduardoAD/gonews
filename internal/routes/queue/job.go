package queue

import (
	"context"

	"github.com/CarlosEduardoAD/go-news/internal/config/task_queue"
	jobcontrollers "github.com/CarlosEduardoAD/go-news/internal/controllers/job_controllers"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

func HandleEmailDelivery(ctx context.Context, t *asynq.Task) error {
	client := task_queue.GenerateAsynqClient()

	controller := jobcontrollers.NewJobController(client)
	err := controller.ExecuteTask(ctx, t)

	if err != nil {
		logrus.Error(shared.GenerateError(err))
		return err
	}

	return nil
}
