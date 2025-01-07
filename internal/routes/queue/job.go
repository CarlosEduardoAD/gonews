package queue

import (
	"log"

	"github.com/CarlosEduardoAD/go-news/internal/config/task_queue"
	jobcontrollers "github.com/CarlosEduardoAD/go-news/internal/controllers/job_controllers"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/gocraft/work"
	"github.com/sirupsen/logrus"
)

func HandleEmailDelivery(work *work.Job) error {
	client := task_queue.GenerateAsynqClient()

	controller := jobcontrollers.NewJobController(client)
	err := controller.ExecuteTask(work)

	if err != nil {
		log.Println("deu bizil: ", err)
		logrus.Error(shared.GenerateError(err))
		return err
	}

	return nil
}
