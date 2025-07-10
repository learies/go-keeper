package config

import (
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - структура конфигурации приложения
type Config struct {
	Env string `yaml:"env" env-default:"local" env:"ENV"`
	Log struct {
		Level string `yaml:"level" env-default:"debug" env:"LOG_LEVEL"`
	} `yaml:"log"`
	Server struct {
		Host string `yaml:"host" env-default:"localhost" env:"SERVER_HOST"`
		Port string `yaml:"port" env-default:"8080" env:"SERVER_PORT"`
	} `yaml:"server"`
	GRPC struct {
		Host string `yaml:"host" env-default:"localhost" env:"GRPC_HOST"`
		Port string `yaml:"port" env-default:"50051" env:"GRPC_PORT"`
	} `yaml:"grpc"`
}

// ParseLogLevel преобразует строку в соответствующий slog.Level
func ParseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// fetchConfigPath - получает путь до конфигурационного файла
// Приоритеты:
// 1. Флаг командной строки -config
// 2. Переменная окружения CONFIG_PATH
// Если ничего не задано - возвращает пустую строку
func fetchConfigPath() string {
	var res string

	// Регистрируем флаг -config для командной строки
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	// Если флаг не задан, проверяем переменную окружения
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

// MustLoadConfig - загружает конфигурацию из файла
// Паникует, если:
// - не указан путь до конфигурационного файла
// - файл не существует
// - произошла ошибка при чтении файла
func MustLoadConfig() *Config {
	// Получаем путь до конфигурационного файла
	path := fetchConfigPath()

	// Проверяем, что путь не пустой
	if path == "" {
		panic("config path is empty")
	}

	// Проверяем, что файл существует
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	// Читаем конфигурацию из файла
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
