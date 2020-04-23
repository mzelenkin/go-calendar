// Пакет предоставляет структурирное журналирование с logrus.
package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

var (
	// Logger сконфигурированный logrus.Logger.
	Logger *logrus.Logger

	// Файл журнала или nil
	logFile *os.File
)

// HTTPStructuredLogger структурированный logrus Logger.
type HTTPStructuredLogger struct {
	Logger *logrus.Logger
}

// NewLogger создает и конфигурирует logrus Logger.
func NewLogger() *logrus.Logger {
	if Logger != nil {
		return Logger
	}

	Logger = logrus.New()

	setOutput()
	setFormatter()
	setLevel()

	// Устанавливаем вывод стандартного логгера в logrus
	log.SetOutput(Logger.Writer())

	return Logger
}

// Exit вызывается при выходе из программы
// Это обычная обертка над logrus.Exit, нужна для корректного закрытия файла журнала
// (вызова хэндлеров, зарегистрированных в RegisterExitHandler)
func Exit(code int) {
	if Logger != nil {
		Logger.Exit(code)
	}
}

// setOutput устанавливает вывод журнала в зависимости от настроек
func setOutput() {
	var err error

	filename := viper.GetString("log.file")
	if filename != "" {
		openFlags := os.O_CREATE | os.O_WRONLY
		if viper.GetBool("log.file_append") {
			openFlags |= os.O_APPEND
		} else {
			openFlags |= os.O_TRUNC
		}

		logFile, err = os.OpenFile(filename, openFlags, 0650)
		if err != nil {
			log.Fatal(err)
		}

		Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
		logrus.SetOutput(io.MultiWriter(os.Stdout, logFile)) // На всякий устанавливаем поток вывода глобально
		logrus.RegisterExitHandler(closeLogFile)             // Закрыватор файла журнала при выходе
	}
}

// setFormatter настраивает формат вывода
// Флаг настроек textlogging определяет в каком формате вудет осуществляться вывод
// При textlogging =  true данные выводятся в текстовом формате, иначе в формате JSON
func setFormatter() {
	if viper.GetBool("log.textlogging") {
		Logger.Formatter = &logrus.TextFormatter{}
	} else {
		Logger.Formatter = &logrus.JSONFormatter{}
	}
}

// setLevel устанавливает уровень событий для которых происходит журналирование
// Настройки определяет параметр level в конфигурации
func setLevel() {
	level := viper.GetString("log.level")
	if level == "" {
		level = "info"
	}
	l, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	Logger.Level = l
}

// closeLogFile закрывает файл журнала при завершении журналирования
func closeLogFile() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, "Logfile closing error: ", err)
		}
	}
}
