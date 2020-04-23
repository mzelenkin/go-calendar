package logging

import (
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

// NewStructuredLogger реализует кастомный структурированный логгер
func NewHTTPStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&HTTPStructuredLogger{logger})
}

// NewHTTPLogEntry создает структурированный логгер и добавляет поля HTTP-запроса
func (l *HTTPStructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	// Если у нас есть ID запроса, пишем его в логи
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	// Определяем что за схема используется (http/https)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Добавляем информацию о запросе
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	//logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	logFields["uri"] = r.RequestURI

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("request started")

	return entry
}

// Вспомогательные методы, используемые приложением для получения логгера
// с уже установленными полями запроса, а также для установки дополнительных полей.

// GetHTTPLogEntry возвращает логгер с предустановленными полями
// в котором уже журналируются основные параметры запроса
func GetHTTPLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

// LogHTTPEntrySetField добавляет поле из запроса в logrus.FieldLogger.
func LogHTTPEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

// LogHTTPEntrySetFields добавляет несколько полей из запроса в logrus.FieldLogger.
func LogHTTPEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
