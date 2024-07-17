package db

import (
	"context"
	"fmt"

	"github.com/alexPavlikov/go-atm/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Функция подключения к базе данных
func Connect(cfg *config.Config) (conn *pgxpool.Pool, err error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.PostgreUser, cfg.PostgrePassword, cfg.PostgresPath, cfg.PostgresPort, cfg.PostgreDatabaseName)
	conn, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
