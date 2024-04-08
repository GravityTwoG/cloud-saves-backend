package email_sender

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
	) error
}

type mailSender struct {
	name              string
	fromEmailAddress  string
	smtpServerAddress string
	smtpAuth          smtp.Auth
}

func NewEmailSender(
	name string,
	fromEmailAddress string,
	fromEmailPassword string,
	smtpAuthAddress string,
	smtpServerAddress string,
) EmailSender {
	return &mailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		smtpServerAddress: smtpServerAddress,
		smtpAuth:          smtp.PlainAuth("", fromEmailAddress, fromEmailPassword, smtpAuthAddress),
	}
}

func (s *mailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.name, s.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	return e.Send(s.smtpServerAddress, s.smtpAuth)
}
