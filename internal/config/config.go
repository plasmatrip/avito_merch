package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Host        string `env:"RUN_ADDRESS"`  //адрес веб-сервера
	Database    string `env:"DATABASE_URI"` //DSN базы данных
	LogLevel    string `env:"LOG_LEVEL"`    //уровень логирования
	TokenSecret string `env:"TOKEN_SECRET"` //секретный ключ для JWT
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	//пытаемся загрузить .env файл
	if err := godotenv.Load(filepath.Dir(ex) + "/.env"); err != nil {
		return nil, errors.New(".env not found")
	}

	// читаем переменные окружения, при ошибке прокидываем ее наверх
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variable: %w", err)
	}

	if _, exist := os.LookupEnv("RUN_ADDRESS"); !exist {
		return nil, errors.New("RUN_ADDRESS not found")
	}

	if _, exist := os.LookupEnv("DATABASE_URI"); !exist {
		return nil, errors.New("DATABASE_URI not found")
	}

	if _, exist := os.LookupEnv("LOG_LEVEL"); !exist {
		return nil, errors.New("LOG_LEVEL not found")
	}

	if _, exist := os.LookupEnv("TOKEN_SECRET"); !exist {
		return nil, errors.New("TOKEN_SECRET not found")
	}

	return cfg, nil
}
