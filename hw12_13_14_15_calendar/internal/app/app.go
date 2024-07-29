package app

import (
	"context"
	"fmt"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	"log"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/server/internalhttp"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/services"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/storage/sql"
)

type App struct {
	config              *config.Config
	logger              logger.Logger
	httpServer          *internalhttp.Server
	grpcServer          *grpc.Server
	storage             storage.Storage
	eventService        services.EventService
	notificationService services.NotificationService
	healthService       services.HealthService
}

func NewApp(config *config.Config) (*App, error) {
	// Инициализация логгера
	logInstance, err := logger.New(config.Logger)
	if err != nil {
		log.Fatalf("on initializing logger, %s", err)
	}

	store, err := initStorage(config.Database, logInstance)
	if err != nil {
		return nil, fmt.Errorf("on initializing storage, %w", err)
	}

	app := &App{
		config:  config,
		logger:  logInstance,
		storage: store,
	}

	// Инициализация сервисов
	app.eventService = services.NewEventService(store.EventRepository())
	app.notificationService = services.NewNotificationService(store.NotificationRepository())
	app.healthService = services.NewHealthService(store)

	// Initialize servers
	app.httpServer = internalhttp.New(
		config.HTTPServer,
		logInstance,
		app.eventService,
		app.notificationService,
		app.healthService,
	)

	grpcServer, err := grpc.New(
		app.eventService,
		app.notificationService,
		app.healthService,
		logInstance,
		config.GRPCServer,
	)
	if err != nil {
		return nil, fmt.Errorf("on initializing gRPC server, %w", err)
	}
	app.grpcServer = grpcServer

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info("Календарь запущен...")
	a.logger.Info(a.config)
	// Запуск HTTP сервера
	go func() {
		if err := a.httpServer.Start(ctx); err != nil {
			a.logger.Error("failed to start HTTP server: " + err.Error())
		}
		a.logger.Info("http сервер запущен")
	}()

	// Запуск gRPC сервера
	go func() {
		if err := a.grpcServer.Start(ctx); err != nil {
			a.logger.Error("failed to start gRPC server: " + err.Error())
		}
	}()

	<-ctx.Done()
	return a.Stop(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	// Остановка HTTP сервера
	if err := a.httpServer.Stop(ctx); err != nil {
		return err
	}

	// Остановка gRPC сервера
	a.grpcServer.Stop(ctx)
	return nil
}

func initStorage(config config.DatabaseConfig, logger logger.Logger) (store storage.Storage, err error) {
	switch config.Storage {
	case "memory":
		store = memorystorage.New()
	case "sql":
		store, err = sqlstorage.New(config, logger)
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
