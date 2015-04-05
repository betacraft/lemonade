package mailer

import (
	"github.com/jordan-wright/email"
	"net/smtp"
)

var auth *smtp.Auth

func Init() {
	tAuth := smtp.PlainAuth("", "lemonades@rainingclouds.com", "lemonades1511", "smtp.gmail.com")
	auth = &tAuth
}

func SendEmail(to []string, content []byte) error {
	return smtp.SendMail("smtp.gmail.com:587", *auth, "lemonades@rainingclouds.com", to, content)
}

func Send(to string, email *email.Email) {
	email.To = []string{to}
	email.Send("smtp.gmail.com:587", *auth)
}

func SendToMany(to []string, email *email.Email) {
	email.To = to
	email.Send("smtp.gmail.com:587", *auth)
}
