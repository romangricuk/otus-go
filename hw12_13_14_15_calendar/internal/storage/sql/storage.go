package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	// "github.com/lib/pq" подключение драйвера postgres.
	_ "github.com/lib/pq"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type SQLStorage struct {
	db               *sql.DB
	eventRepo        *EventRepo
	notificationRepo *NotificationRepo
	logger           logger.Logger
}

func New(cfg config.DatabaseConfig, logger logger.Logger) (*SQLStorage, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &SQLStorage{
		db:               db,
		eventRepo:        NewEventRepo(db, logger),
		notificationRepo: NewNotificationRepo(db, logger),
		logger:           logger,
	}, nil
}

func (s *SQLStorage) Connect(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

func (s *SQLStorage) EventRepository() storage.EventRepository {
	return s.eventRepo
}

func (s *SQLStorage) NotificationRepository() storage.NotificationRepository {
	return s.notificationRepo
}

func (s *SQLStorage) HealthCheck(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
