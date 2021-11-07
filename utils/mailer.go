package utils

import (
	"net/smtp"
	"os"
)

// Mailer has just one method. To send a simple mail
type Mailer interface {
	SendSimpleMail(to string, body []byte) error
}

// NewSimpleMail is a constructor for Mailer interface. Returns a simpleMail struct
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

// SendSimpleMail utilizes the golang smtp library to send a simple mail
func (s simpleMail) SendSimpleMail(to string, body []byte) error {
	auth := smtp.PlainAuth("Stud Pal", s.from, s.password, s.smtpHost)
	return smtp.SendMail(s.smtpHost + ":" + s.smtpPort, auth, s.from, []string{to}, body)
}
