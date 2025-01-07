package queue

import (
	"log"

	jobcontrollers "github.com/CarlosEduardoAD/go-news/internal/controllers/job_controllers"
	"github.com/CarlosEduardoAD/go-news/internal/shared"
	"github.com/gocraft/work"
	"github.com/sirupsen/logrus"
)

func HandleEmailDelivery(work *work.Job) error {
	controller := jobcontrollers.NewJobController()
	err := controller.ExecuteTask(work)

	if err != nil {
		log.Println(err)
		logrus.Error(shared.GenerateError(err))
		return err
	}

	return nil
}
