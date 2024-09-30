package config

import (
	"fmt"
	"os"
	"strings"

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
	viper.SetDefault("email.insecureSkipVerify", true)

	// Настройка замены переменных окружения
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Устанавливаем тип конфигурационного файла на YAML
	viper.SetConfigType("yaml")

	// Чтение конфигурации из файла
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Замена переменных окружения в файле конфигурации
	configContentStr := os.ExpandEnv(string(configContent))

	// Устанавливаем обработанный контент в Viper
	if err := viper.ReadConfig(strings.NewReader(configContentStr)); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
