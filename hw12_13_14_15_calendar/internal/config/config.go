package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServerConfig
	GRPCServer GRPCServerConfig
	Database   DatabaseConfig
	Logger     LoggerConfig
	RabbitMQ   RabbitMQConfig
	Sender     SenderConfig
	Scheduler  SchedulerConfig
	Email      EmailConfig
}

type HTTPServerConfig struct {
	Address string
}

type GRPCServerConfig struct {
	Address string
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     int
	Storage  string
}

type LoggerConfig struct {
	Level            string
	Encoding         string
	OutputPaths      []string
	ErrorOutputPaths []string
}

type RabbitMQConfig struct {
	URL       string
	QueueName string
}

type EmailConfig struct {
	SMTPServer         string
	SMTPPort           int
	Username           string
	Password           string
	From               string
	UseTLS             bool
	InsecureSkipVerify bool
}

type SenderConfig struct {
	Interval int // Интервал для проверки очереди RabbitMQ в секундах
}

type SchedulerConfig struct {
	Interval int // Интервал выполнения задач в секундах
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	// Устанавливаем значения по умолчанию
	viper.SetDefault("httpserver.address", "0.0.0.0:8080")
	viper.SetDefault("grpcserver.address", "0.0.0.0:9090")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "calendar")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.storage", "sql")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.encoding", "json")
	viper.SetDefault("logger.outputPaths", []string{"stdout"})
	viper.SetDefault("logger.errorOutputPaths", []string{"stderr"})
	viper.SetDefault("rabbitmq.url", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("rabbitmq.queueName", "calendar_queue")
	viper.SetDefault("sender.interval", 10)
	viper.SetDefault("scheduler.interval", 10)
	viper.SetDefault("email.useTLS", false)

	// Загрузка переменных окружения
	if err := bindEnvVariables(); err != nil {
		return nil, err
	}

	// Чтение конфигурации
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func bindEnvVariables() error {
	envBindings := map[string]string{
		"DB_USER":             "database.user",
		"DB_PASSWORD":         "database.password",
		"DB_NAME":             "database.name",
		"DB_HOST":             "database.host",
		"DB_PORT":             "database.port",
		"RABBITMQ_URL":        "rabbitmq.url",
		"RABBITMQ_QUEUE_NAME": "rabbitmq.queueName",
		"SENDER_INTERVAL":     "sender.interval",
		"SCHEDULER_INTERVAL":  "scheduler.interval",
		"EMAIL_SMTP_SERVER":   "email.smtpServer",
		"EMAIL_SMTP_PORT":     "email.smtpPort",
	}

	for envVar, configKey := range envBindings {
		if err := viper.BindEnv(configKey, envVar); err != nil {
			return fmt.Errorf("on bind env var: %w", err)
		}
	}

	return nil
}
