package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func GenerateLogrus() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetOutput(file)

	return logger
}
