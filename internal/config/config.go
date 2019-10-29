package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/pkg/errors"
)

// LoadConfig загружает конфигурацию и возвращает ее. В случае ошибки будет возвращена err
func LoadConfig(configFilename string) (*Config, error) {

	// Для загрузки конфигурации используем confita
	loader := confita.NewLoader(
		file.NewBackend(configFilename),
	)

	cfg := createDefaultConfig()

	if err := loader.Load(context.Background(), cfg); nil != err {
		// в случае ошибки возвращаем ошибку, путь вызвавший нас метод сам примет решение что делать
		return nil, errors.Errorf("Failed to load app configuration: %+v", err)
	}

	return cfg, nil
}

// Создает конфигурацию по умолчанию
func createDefaultConfig() *Config {
	return &Config{
		Listen: Listen{9001, "127.0.0.1"},
		Log:LoggingConfig{Level:"debug"},
	}
}