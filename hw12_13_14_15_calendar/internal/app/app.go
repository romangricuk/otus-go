package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

type Application interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App struct {
	config              *config.Config
	logger              logger.Logger
	server              *internalhttp.Server
	storage             storage.Storage
	eventService        services.EventService
	notificationService services.NotificationService
	healthService       services.HealthService
}

func NewApp(config *config.Config) (*App, error) {
	// Initialize logger
	logInstance, err := logger.New(config.Logger)
	if err != nil {
		log.Fatalf("on initializing logger, %s", err)
	}

	store, err := initStorage(config.Database)
	if err != nil {
		logInstance.Fatalf("on initializing storage, %s", err)
	}

	app := &App{
		config:  config,
		logger:  logInstance,
		storage: store,
	}

	// Initialize services
	app.eventService = services.NewEventService(store.EventRepository())
	app.notificationService = services.NewNotificationService(store.NotificationRepository())
	app.healthService = services.NewHealthService(store)

	// Initialize server
	app.server = internalhttp.New(config.Server, logInstance, app.eventService, app.notificationService, app.healthService)

	return app, nil
}

func initStorage(config config.DatabaseConfig) (store storage.Storage, err error) {
	switch config.Storage {
	case "memory":
		store = memorystorage.New()
	case "sql":
		store, err = sqlstorage.New(config)
		if err != nil {
			return nil, fmt.Errorf("on initializing SQL storage, %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown storage type: %s", config.Storage)
	}

	// Connect to storage
	if err := store.Connect(context.Background()); err != nil {
		return nil, fmt.Errorf("on connecting to storage, %w", err)
	}

	return store, nil
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info("calendar is running...")
	return a.server.Start(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	errServer := a.server.Stop(ctx)
	if errServer != nil {
		errServer = fmt.Errorf("on stopping server, %w", errServer)
	}

	errStorage := a.storage.Close()
	if errStorage != nil {
		errStorage = fmt.Errorf("on closing storage, %w", errStorage)
	}

	return errors.Join(errServer, errStorage)
}
