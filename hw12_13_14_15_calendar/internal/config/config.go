package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServerConfig
	GRPCServer GRPCServerConfig
	Database   DatabaseConfig
	Logger     LoggerConfig
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

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	// Устанавливаем значения по умолчанию
	viper.SetDefault("server.address", "0.0.0.0:8080")
	viper.SetDefault("grpc.address", "0.0.0.0:9090")
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

	// Загрузка переменных окружения
	viper.AutomaticEnv()

	// Чтение конфигурации
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Замена переменных окружения в конфигурации
	replaceEnvVariables(&config)

	return &config, nil
}

func replaceEnvVariables(cfg *Config) {
	cfg.Database.User = getEnv("DB_USER", cfg.Database.User)
	cfg.Database.Password = getEnv("DB_PASSWORD", cfg.Database.Password)
	cfg.Database.Name = getEnv("DB_NAME", cfg.Database.Name)
	cfg.Database.Host = getEnv("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvAsInt("DB_PORT", cfg.Database.Port)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
