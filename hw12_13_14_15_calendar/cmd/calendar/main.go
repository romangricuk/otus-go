package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
)

func main() {
	// Парсинг флагов командной строки
	fmt.Println("app started")

	configPath := flag.String("config", "configs/config.toml", "path to the config file")
	migrationsPath := flag.String("migrations", "migrations", "path to the migrations directory")
	command := flag.String("command", "run", "command to execute: run, migrate_up, migrate_down")
	flag.Parse()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Выбор команды для выполнения
	switch *command {
	case "run":
		fmt.Println("command is run")

		runApplication(cfg)
	case "migrate_up":
		fmt.Println("command is migrate_up")

		runMigrations(cfg, *migrationsPath, "up")
	case "migrate_down":
		fmt.Println("command is migrate_down")

		runMigrations(cfg, *migrationsPath, "down")
	default:
		log.Fatalf("Unknown command: %s. Use 'run', 'migrate_up' or 'migrate_down'", *command)
	}
}

func runApplication(cfg *config.Config) {
	// Инициализация приложения
	fmt.Println("runApplication")

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Error initializing application, %s", err)
	}

	// Контекст и корректное завершение работы
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := application.Stop(ctx); err != nil {
			log.Printf("failed to stop application: %s", err)
		}
	}()

	if err := application.Start(ctx); err != nil {
		log.Printf("failed to start application: %s", err)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func runMigrations(cfg *config.Config, migrationsPath string, direction string) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	migrationPath := fmt.Sprintf("file://%s", filepath.Clean(migrationsPath))

	m, err := migrate.New(
		migrationPath,
		dsn)
	if err != nil {
		log.Fatalf("Error initializing migration, %s", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Error applying migrations, %s", err)
		}
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Error reverting migrations, %s", err)
		}
	default:
		log.Fatalf("Unknown migration direction: %s. Use 'up' or 'down'", direction)
	}

	fmt.Println("Migrations applied successfully.")
}
