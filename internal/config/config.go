package config

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Config struct {
	Env                 string        `mapstructure:"ENV"`
	Timeout             time.Duration `mapstructure:"TIMEOUT"`
	LogLevel            string        `mapstructure:"LOGLEVEL"`
	Path                string        `mapstructure:"SERVER_PATH"`
	Port                int           `mapstructure:"SERVER_PORT"`
	PostgresPath        string        `mapstructure:"POSTGRES_PATH"`
	PostgresPort        int           `mapstructure:"POSTGRES_PORT"`
	PostgreUser         string        `mapstructure:"POSTGRES_USER"`
	PostgrePassword     string        `mapstructure:"POSTGRES_PASSWORD"`
	PostgreDatabaseName string        `mapstructure:"POSTGRES_DATABASE"`
}

func (c *Config) ServerToString() string {
	return fmt.Sprintf("%s:%d", c.Path, c.Port)
}

// Функция загрузки конфига из файла
func Load() (*Config, error) {

	path, file := fetchConfigPath()

	if path == "" || file == "" {
		return &Config{}, errors.New("fetch config path error")
	}

	var cfg Config

	cfg, err := initViper(path, file, cfg)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to init viper: %w", err)
	}

	slog.Info("config file load successfully", "log_level", cfg.LogLevel)

	return &cfg, nil
}

// Функция получение пути к файлу конфига
func fetchConfigPath() (path, file string) {

	flag.StringVar(&path, "config_path", "", "config file path")
	flag.StringVar(&file, "config_file", "", "config file name")
	flag.Parse()

	return path, file
}

// Функция инициализация viper для считывания конфиг файла в структуру
func initViper(path string, file string, cfg Config) (Config, error) {
	viper.SetConfigName(file)
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Функция для инициализации слоя логирования
func SetupLogger(logLevel string) {
	var log *slog.Logger

	switch logLevel {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	slog.SetDefault(log)
}
