package service

import (
	"context"
	"fmt"
	"log/slog"

	models "github.com/alexPavlikov/go-atm/internal/domain"
	postgres "github.com/alexPavlikov/go-atm/internal/server/db"
)

type Services struct {
	postgres *postgres.Postgres
}

func New(postgres *postgres.Postgres) *Services {
	return &Services{
		postgres: postgres,
	}
}

func (s *Services) AddAccount(ctx context.Context, account models.UserAccount) (id int, err error) {

	id, err = s.postgres.InsertAccount(ctx, account)
	if err != nil {
		slog.Error("failed to insert account to postgres", "error", err)
		return -1, fmt.Errorf("failed to insert account to postgres: %w", err)
	}

	return id, nil
}

func (s *Services) UpDeposite(ctx context.Context, amount float64, id int) error {
	if err := s.postgres.Deposite(ctx, amount, id); err != nil {
		slog.Error("failed to deposite", "error", err)
		return fmt.Errorf("failed to up deposite: %w", err)
	}

	return nil
}

func (s *Services) Withdraw(ctx context.Context, amount float64, id int) error {
	if err := s.postgres.Withdraw(ctx, amount, id); err != nil {
		slog.Error("failed to withdraw", "error", err)
		return fmt.Errorf("failed to up deposite: %w", err)
	}
	return nil
}

func (s *Services) GetAccount(ctx context.Context, id int) (models.UserAccount, error) {
	acc, err := s.postgres.SelectAccount(ctx, id)
	if err != nil {
		slog.Error("failed to get user account", "error", err)
		return models.UserAccount{}, fmt.Errorf("failed to get user account: %w", err)
	}
	return acc, nil
}
