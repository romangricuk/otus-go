package services

import (
	"fmt"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/email"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type SenderService struct {
	emailClient email.Client
	logger      logger.Logger
}

func NewSenderService(emailClient email.Client, logger logger.Logger) *SenderService {
	return &SenderService{
		emailClient: emailClient,
		logger:      logger,
	}
}

func (s *SenderService) ProcessNotification(notification dto.NotificationData) error {
	// Логируем полученное уведомление
	s.logger.Info("Received notification: %+v", notification)

	// Отправляем email
	if err := s.sendEmailNotification(notification); err != nil {
		return fmt.Errorf("send email notification failed: %w", err)
	}

	return nil
}

func (s *SenderService) sendEmailNotification(notification dto.NotificationData) error {
	// Формируем email
	to := "example@example.com" // Тут могла бы быть обработка получения EMAIL от пользователя
	subject := fmt.Sprintf("Calendar Notification #%s", notification.ID)
	body := notification.Message

	// Отправляем email
	if err := s.emailClient.SendEmail(to, subject, body); err != nil {
		s.logger.Error("on send email: %v", err)
		return err
	}

	return nil
}
