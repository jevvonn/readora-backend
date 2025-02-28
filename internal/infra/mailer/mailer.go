package mailer

import (
	"github.com/jevvonn/readora-backend/config"
	"gopkg.in/gomail.v2"
)

type MailerItf interface {
	Send(to []string, subject, body string) error
}

type Mailer struct {
	Dialer *gomail.Dialer
	Conf   config.Config
}

func New() MailerItf {
	conf := config.Load()

	d := gomail.NewDialer(
		conf.SMTPHost,
		conf.SMTPPort,
		conf.SMTPUsername,
		conf.SMTPPassword,
	)

	return &Mailer{
		Dialer: d,
		Conf:   conf,
	}
}

func (m *Mailer) Send(to []string, subject, body string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "Readora <"+m.Conf.SMTPEmail+">")
	mail.SetHeader("To", to...)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", body)

	return m.Dialer.DialAndSend(mail)
}
