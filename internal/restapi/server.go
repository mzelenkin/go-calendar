package restapi

import (
	"context"
	"github.com/mzelenkin/go-calendar/internal/logging"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Server struct {
	*http.Server
}

// NewServer конструктор REST API сервера
func NewServer() (*Server, error) {
	log.Println("configuring server...")
	api, err := New(viper.GetBool("http.enable_cors"))
	if err != nil {
		return nil, err
	}

	// Получаем настройки из viper
	addr := viper.GetString("http.listen")

	srv := http.Server{
		Addr:    addr,
		Handler: api,
	}

	return &Server{&srv}, nil
}

// Start запускает сервер, а также корректно отрабатывает завершение его работы
func (srv *Server) Start() {
	logger := logging.NewLogger()
	logger.Info("Starting server")

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	logger.Infof("Server started. Listening on %s\n", srv.Addr)

	// Обработка сигнала завершения операционной системы SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	logger.Warn("Shutting down server... Reason:", sig)

	// Логика завершения
	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	logger.Info("Server gracefully stopped")
}
