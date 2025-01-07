package jobmodel

import (
	"fmt"
	"testing"
	"time"

	"github.com/CarlosEduardoAD/go-news/internal/utils"
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

	err := sendEmailJob.AddAndEnqueueTask()

	assert.NotNil(t, err)
}

func TestValidAddAndEnqueueTask(t *testing.T) {

	sendEmailJob := NewSendEmailJob(utils.GenerateRandomString(8), fmt.Sprintf("%s@test.com", utils.GenerateRandomString(8)), int64(time.Now().Second()), "send_email")

	err := sendEmailJob.AddAndEnqueueTask()

	assert.Nil(t, err)
}
