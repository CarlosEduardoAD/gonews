package utils

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"
)

func GenerateRandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)

}

func ReturnNextMonday() int64 {
	weekDays := 7
	now := time.Now().Local()
	nowDay := int(now.Weekday())

	daysUntilNextMonday := (8 - nowDay) % weekDays

	if daysUntilNextMonday == 0 {
		daysUntilNextMonday = 7 // Se hoje é segunda, vamos para a próxima
	}

	nextMonday := now.AddDate(0, 0, daysUntilNextMonday)
	nextMonday = time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 8, 0, 0, 0, nextMonday.Location())

	return int64(nextMonday.Unix())
}

type EmailData struct {
	Name            string
	ConfirmLink     string
	NewsTitle       string
	NewsDescription string
	NewsLink        string
	UnsubscribeLink string
}

func LoadTemplate(templatePath string, data EmailData) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
