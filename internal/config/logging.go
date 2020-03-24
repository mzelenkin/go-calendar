package config

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
)

// LoggingConfig структура с настройками логгирования
type LoggingConfig struct {
	Level    string `config:"level"`
	Filename string `config:"filename"`
}

// ConfigureLogging настраивает журналирование и возвращает готовый объект логгера
func ConfigureLogging(config *LoggingConfig) (*logrus.Logger, error) {
	if config.Filename != "" {
		f, errOpen := os.OpenFile(config.Filename, os.O_RDWR|os.O_APPEND, 0660)
		if errOpen != nil {
			return nil, errOpen
		}
		logrus.SetOutput(bufio.NewWriter(f))
	}

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logrus.StandardLogger(), nil
}
