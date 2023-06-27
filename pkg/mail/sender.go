package mail

//TODO add aws ses sender

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

var (
	smtpHost       = "smtp.gmail.com"
	smtpServerAddr = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddr     string
	fromEmailPassword string
}

func NewGmailSender(
	name string,
	fromEmailAddr string,
	fromEmailPassword string,
) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddr:     fromEmailAddr,
		fromEmailPassword: fromEmailPassword,
	}
}

func (s *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = s.name + "<" + s.fromEmailAddr + ">"
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.Subject = subject
	e.HTML = []byte(content)

	for _, f := range attachFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", s.fromEmailAddr, s.fromEmailPassword, smtpHost)
	return e.Send(smtpServerAddr, smtpAuth)
}
