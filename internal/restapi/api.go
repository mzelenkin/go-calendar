package restapi

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/mzelenkin/go-calendar/internal/logging"
	"net/http"
	"time"
)

// New конструктор HTTP API на базе Chi.
// Он создает и настраивает необходимые компоненты для работы API
func New(enableCORS bool) (*chi.Mux, error) {
	logger := logging.NewLogger()
	r := chi.NewRouter()

	r.Use(middleware.Recoverer) // Восстановление после паники
	// r.Use(middleware.RequestID) // Присвоение запросу уникального ID
	// r.Use(middleware.RealIP)	// Необходим если наш сервис находится за reverse proxy вроде NGINX
	r.Use(middleware.Timeout(15 * time.Second)) // Таймаут соединения с использованием контекста

	r.Use(logging.NewHTTPStructuredLogger(logger))       // Добавляем логгер запросов
	r.Use(render.SetContentType(render.ContentTypeJSON)) // Устанавливаем тип контента application/json

	// Если включена поддержка Cross-Origin Request Sharing (CORS), используем CORS middleware из пакета chi
	// Требуется браузеру для ослабления правила "одного источника", когда домен запроса не совпадает с доменом API
	if enableCORS {
		r.Use(corsConfig().Handler)
	}

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		log := logging.GetHTTPLogEntry(r)
		log.Info("Requested 'hello' action")

		name := r.URL.Query().Get("name")
		_, err := w.Write([]byte("hello " + name))
		if err != nil {
			log.Error("Response writing error: ", err.Error())
		}
	})

	return r, nil
}

func corsConfig() *cors.Cors {
	// Базовые настройки CORS
	return cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Раскомментировать, если нужно указать конкретные хосты
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                 // Разрешенные методы
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Разрашенные заголовки
		ExposedHeaders:   []string{"Link"},                                                    // Заголовки, которые может читать JS
		AllowCredentials: true,
		MaxAge:           86400, // Максимальное время, на которое предзапрос (CORS preflight) может быть закэширован
	})
}
