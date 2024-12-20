package shared

import (
	"crypto/tls"
	"errors"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Dialer *gomail.Dialer
}

func GenerateEmailSender(host string, port int, user, password string) *EmailSender {
	err := validate(host, port, user, password)

	if err != nil {
		panic(err)
	}

	d := gomail.NewDialer(host, port, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &EmailSender{
		Dialer: d,
	}
}

func validate(host string, port int, user, password string) error {
	if host == "" {
		return errors.New("host field is required")
	}

	if port == 0 {
		return errors.New("port field is required")
	}

	if user == "" {
		return errors.New("user field is required")
	}

	if password == "" {
		return errors.New("password field is required")
	}

	return nil
}

func (es *EmailSender) validateEmail(to, subject, body string) error {
	if to == "" {
		return errors.New("to field is required")
	}

	if subject == "" {
		return errors.New("subject field is required")
	}

	if body == "" {
		return errors.New("body field is required")
	}

	return nil
}

func (es *EmailSender) SendEmail(to string, subject string, body string) error {
	err := es.validateEmail(to, subject, body)

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "karl.devcontato@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := es.Dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
