package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/romangricuk/otus-go/hw12_13_14_15_calendar/internal/config"
)

func main() {
	// Инициализация приложения
	fmt.Println("run sender application")

	configPath := flag.String("config", "configs/config.toml", "path to the config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	application, err := app.NewSchedulerApp(cfg)
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
