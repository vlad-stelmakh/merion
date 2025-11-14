package main

import (
	"errors"
	"fmt"
)

//go:generate mockery --name EmailSender
type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type NotificationService struct {
	emailSender EmailSender
}

func NewNotificationService(emailSender EmailSender) *NotificationService {
	return &NotificationService{
		emailSender: emailSender,
	}
}

func (s *NotificationService) NotifyUser(user *User, message string) error {
	if !user.IsValidEmail() {
		return errors.New("некорректный email пользователя")
	}

	subject := "Уведомление"
	body := fmt.Sprintf("Привет, %s!\n\n%s", user.Name, message)

	return s.emailSender.SendEmail(user.Email, subject, body)
}
