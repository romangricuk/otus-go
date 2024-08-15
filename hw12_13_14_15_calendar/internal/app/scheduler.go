package app

import (
	"context"
	"fmt"
	"time"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	config              *config.Config
	logger              logger.Logger
	notificationService services.NotificationService
	rabbitClient        rabbitmq.Client
	storage             storage.Storage
}

func NewSchedulerApp(cfg *config.Config) (*Scheduler, error) {
	// Инициализация логгера
	logInstance, err := logger.New(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("on initializing logger, %w", err)
	}

	// Инициализация RabbitMQ клиента
	rabbitClient, err := rabbitmq.NewClient(cfg.RabbitMQ, logInstance)
	if err != nil {
		return nil, fmt.Errorf("on initializing RabbitMQ client, %w", err)
	}

	// Инициализация хранилища
	store, err := initStorage(cfg.Database, logInstance)
	if err != nil {
		return nil, fmt.Errorf("on initializing storage, %w", err)
	}

	// Инициализация сервиса уведомлений
	notificationService := services.NewNotificationService(store)

	return &Scheduler{
		config:              cfg,
		logger:              logInstance,
		notificationService: notificationService,
		rabbitClient:        rabbitClient,
		storage:             store,
	}, nil
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Info("Scheduler started")

	// Подключение к хранилищу
	if err := s.storage.Connect(ctx); err != nil {
		return fmt.Errorf("on connecting to storage, %w", err)
	}

	// Подключение к RabbitMQ
	if err := s.rabbitClient.Connect(); err != nil {
		return fmt.Errorf("on connecting to rabbitMQ, %w", err)
	}

	ticker := time.NewTicker(time.Duration(s.config.Scheduler.Interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Scheduler stopped")
			return nil
		case <-ticker.C:
			s.processNotifications(ctx)
			s.cleanupOldNotifications(ctx)
		}
	}
}

func (s *Scheduler) Stop(_ context.Context) error {
	s.logger.Info("Stopping Scheduler")

	if err := s.storage.Close(); err != nil {
		return fmt.Errorf("on close storage connection: %w", err)
	}

	if err := s.rabbitClient.Close(); err != nil {
		return fmt.Errorf("on close rabbit client: %w", err)
	}

	return nil
}

func (s *Scheduler) processNotifications(ctx context.Context) {
	// Получаем уведомления, которые необходимо отправить
	notifications, err := s.notificationService.ListNotifications(ctx, time.Now().Add(-time.Hour*24), time.Now())
	if err != nil {
		s.logger.Errorf("Error listing notifications: %v", err)
		return
	}

	for _, notification := range notifications {
		// Публикуем уведомление в RabbitMQ
		err = s.rabbitClient.SendNotification(notification)
		if err == nil {
			s.logger.Infof("Notification %s published", notification.ID)
		} else {
			s.logger.Errorf("Error publishing notification: %v", err)
		}
	}
}

func (s *Scheduler) cleanupOldNotifications(ctx context.Context) {
	err := s.notificationService.DeleteSentNotifications(ctx)
	if err != nil {
		s.logger.Errorf("Error listing notifications for cleanup: %v", err)
		return
	}
}
