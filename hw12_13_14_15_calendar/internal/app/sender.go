package app

import (
	"context"
	"fmt"
	"log"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/dto"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/email"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
)

type SenderApp struct {
	config       config.SenderConfig
	logger       logger.Logger
	rabbitClient rabbitmq.Client
	service      *services.SenderService
}

func NewSenderApp(cfg *config.Config) (*SenderApp, error) {
	logInstance, err := logger.New(cfg.Logger)
	if err != nil {
		log.Fatalf("on initializing logger, %s", err)
	}

	rabbitClient, err := rabbitmq.NewClient(cfg.RabbitMQ, logInstance)
	if err != nil {
		return nil, err
	}

	emailClient := email.NewSMTPClient(&cfg.Email)
	service := services.NewSenderService(emailClient, logInstance)

	return &SenderApp{
		config:       cfg.Sender,
		logger:       logInstance,
		rabbitClient: rabbitClient,
		service:      service,
	}, nil
}

func (a *SenderApp) Start(ctx context.Context) error {
	a.logger.Info("Starting SenderApp")

	// Подключаемся к RabbitMQ
	if err := a.rabbitClient.Connect(); err != nil {
		return err
	}

	// Запускаем обработку сообщений из очереди
	go func() {
		if err := a.runMessageProcessor(ctx); err != nil {
			a.logger.Error("Failed to process messages: %v", err)
			// Handle errors as needed
		}
	}()

	<-ctx.Done()
	return a.Stop(ctx)
}

func (a *SenderApp) runMessageProcessor(ctx context.Context) error {
	notificationChannel, err := a.rabbitClient.ReceiveNotifications(ctx)
	if err != nil {
		return fmt.Errorf("failed to receive notifications: %w", err)
	}

	for {
		select {
		case notification, ok := <-notificationChannel:
			if !ok {
				return nil // Channel closed, exit
			}
			a.handleNotification(notification)
		case <-ctx.Done():
			return nil // Context canceled, exit
		}
	}
}

func (a *SenderApp) handleNotification(notification dto.NotificationData) {
	// Обработка уведомления
	if err := a.service.ProcessNotification(notification); err != nil {
		a.logger.Error("Failed to process notification: %v", err)
	}
}

func (a *SenderApp) Stop(_ context.Context) error {
	a.logger.Info("Stopping SenderApp")
	err := a.rabbitClient.Close()
	if err != nil {
		return fmt.Errorf("on close rabbit client: %w", err)
	}
	return nil
}
