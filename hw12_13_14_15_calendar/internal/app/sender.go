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
	config              config.SenderConfig
	logger              logger.Logger
	rabbitClient        rabbitmq.Client
	senderService       *services.SenderService
	notificationService services.NotificationService
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

	// Инициализация хранилища
	store, err := initStorage(cfg.Database, logInstance)
	if err != nil {
		return nil, fmt.Errorf("on initializing storage, %w", err)
	}

	// Инициализация сервиса уведомлений
	notificationService := services.NewNotificationService(store)

	return &SenderApp{
		config:              cfg.Sender,
		logger:              logInstance,
		rabbitClient:        rabbitClient,
		senderService:       service,
		notificationService: notificationService,
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
			a.handleNotification(ctx, notification)
		case <-ctx.Done():
			return nil // Context canceled, exit
		}
	}
}

func (a *SenderApp) handleNotification(ctx context.Context, notification dto.NotificationData) {
	// Обработка уведомления
	if err := a.senderService.ProcessNotification(notification); err != nil {
		a.logger.Errorf("Failed to process notification: %v", err)
	}

	notification.Sent = dto.NotificationSent
	err := a.notificationService.UpdateNotification(ctx, notification.ID, notification)
	if err != nil {
		a.logger.Errorf("error updating notification: %w", err)
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
