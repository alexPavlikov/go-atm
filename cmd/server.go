package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexPavlikov/go-atm/internal/config"
	"github.com/alexPavlikov/go-atm/internal/db"
	router "github.com/alexPavlikov/go-atm/internal/server"
	postgres "github.com/alexPavlikov/go-atm/internal/server/db"
	"github.com/alexPavlikov/go-atm/internal/server/locations"
	"github.com/alexPavlikov/go-atm/internal/server/service"
)

// Функция инициализации и запуска сервера
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	srv, err := ServerLoad(cfg)
	if err != nil {
		return err
	}

	// load http server
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Path, cfg.Port), srv); err != nil {
		slog.Error("listen and serve server error", "error", err.Error())
		return err
	}

	return nil
}

func ServerLoad(cfg *config.Config) (http.Handler, error) {
	// setup logger
	config.SetupLogger(cfg.LogLevel)

	slog.Info("starting application", "server config", fmt.Sprintf("%s:%d", cfg.Path, cfg.Port))

	// init handler request
	slog.Info("initialization driver handlers")

	conn, err := db.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed connect to database: %w", err)
	}

	repository := postgres.New(conn)

	service := service.New(repository)

	handlers := locations.New(*service)

	serverBuilder := router.New(handlers)

	srv := serverBuilder.Build()

	return srv, nil
}
