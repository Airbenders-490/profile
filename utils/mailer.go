package utils

import (
	"net/smtp"
	"os"
)

type Mailer interface {
	SendSimpleMail(to string, body []byte) error
}

func NewSimpleMail() Mailer {
	return simpleMail{
		from:     os.Getenv("EMAIL"),
		password: os.Getenv("PASSWORD"),
		smtpHost: "smtp.gmail.com",
		smtpPort: "587",
	}
}

type simpleMail struct {
	from string
	password string
	// smtp.gmail.com
	smtpHost string
	// 587
	smtpPort string
}

func (s simpleMail) SendSimpleMail(to string, body []byte) error {
	auth := smtp.PlainAuth("Stud Pal", s.from, s.password, s.smtpHost)
	return smtp.SendMail(s.smtpHost + ":" + s.smtpPort, auth, s.from, []string{to}, body)
}
