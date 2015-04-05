package mailer

import (
	"github.com/jordan-wright/email"
	"testing"
)

func TestSendToMany(t *testing.T) {
	Init()
	mail := email.NewEmail()
	mail.From = "lemonades@rainingclouds.com"
	mail.Subject = "Test Subject"
	mail.Text = []byte("Test")
	SendToMany([]string{"akshay@rainingclouds.com"}, mail)
}
