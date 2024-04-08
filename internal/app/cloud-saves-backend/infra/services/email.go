package services

import (
	"cloud-saves-backend/internal/app/cloud-saves-backend/domain/user"
	email_sender "cloud-saves-backend/internal/app/cloud-saves-backend/infra/email-sender"
	"context"
	"fmt"
)

type EmailService interface {
	SendPasswordResetEmail(
		ctx context.Context,
		user *user.User,
		token string,
	) error
}

type emailService struct {
	mailer     email_sender.EmailSender
	apiBaseURL string
}

func NewEmail(mailer email_sender.EmailSender, apiBaseURL string) EmailService {
	return &emailService{
		mailer:     mailer,
		apiBaseURL: apiBaseURL,
	}
}

func (s *emailService) SendPasswordResetEmail(
	ctx context.Context,
	user *user.User,
	token string,
) error {
	url := fmt.Sprintf(
		`cloud-saves://reset-password?token=%s`,
		token,
	)
	content := fmt.Sprintf(
		`Hello <b>%s</b>, here is your password reset link: 
		<a href="%s" href="%s/redirect/?redirect-to=%s" href="%s">link</a>
		<br>
		<br>
		%s`,
		user.GetUsername(),
		url,
		s.apiBaseURL,
		url,
		url,
		url,
	)

	return s.mailer.SendEmail(
		"Password Reset",
		content,
		[]string{user.GetEmail()},
		[]string{},
		[]string{},
	)
}
