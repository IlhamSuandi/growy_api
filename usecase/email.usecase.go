package usecase

import (
	"github.com/ilhamSuandi/business_assistant/config"
	"github.com/ilhamSuandi/business_assistant/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailUsecase interface {
	SendEmail(to, subject, body string) error
	SendEmailAsync(to, subject, body string)
}

type emailUsecase struct {
	Logger *logrus.Logger
	Dialer *gomail.Dialer
}

func NewEmailUsecase() EmailUsecase {
	return &emailUsecase{
		Logger: utils.Log,
		Dialer: gomail.NewDialer(
			config.SMTPHOST,
			config.SMTPPORT,
			config.SMTPUSERNAME,
			config.SMTPPASSWORD,
		),
	}
}

func (eu *emailUsecase) SendEmail(to, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.EMAIL_FROM)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	eu.Logger.Infof("sending email to %s", to)

	// Reuse connection
	sender, err := eu.Dialer.Dial()
	if err != nil {
		eu.Logger.Errorf("Failed to establish SMTP connection: %v", err)
		return err
	}
	defer sender.Close()

	if err := gomail.Send(sender, mailer); err != nil {
		eu.Logger.Errorf("Failed to send email: %v", err)
		return err
	}

	eu.Logger.Info("Email sent successfully")
	return nil
}

func (eu *emailUsecase) SendEmailAsync(to, subject, body string) {
	go func() {
		if err := eu.SendEmail(to, subject, body); err != nil {
			eu.Logger.Errorf("Failed to send email asynchronously: %v", err)
		}
	}()
}
