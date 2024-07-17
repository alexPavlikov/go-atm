package postgres

import (
	"context"

	models "github.com/alexPavlikov/go-atm/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

func New(DB *pgxpool.Pool) *Postgres {
	return &Postgres{
		DB: DB,
	}
}

func (p *Postgres) SelectAccount(ctx context.Context, id int) (models.UserAccount, error) {
	query := `SELECT id, name, balance FROM public."Account" WHERE id = $1`

	row := p.DB.QueryRow(ctx, query, id)

	var acc models.UserAccount

	if err := row.Scan(&acc.ID, &acc.Name, &acc.Balance); err != nil {
		return models.UserAccount{}, err
	}

	return acc, nil
}

func (p *Postgres) InsertAccount(ctx context.Context, account models.UserAccount) (id int, err error) {

	query := `
	INSERT INTO public."Account" (name, balance) VALUES ($1, $2) RETURNING id
	`

	row := p.DB.QueryRow(ctx, query, account.Name, int(account.Balance*100))

	if err = row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (p *Postgres) Deposite(ctx context.Context, amount float64, id int) error { //transaction

	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `UPDATE public."Account" SET balance = balance + $1 WHERE id = $2 RETURNING balance`

	row := tx.QueryRow(ctx, query, int(amount*100), id)

	var balace int

	if err := row.Scan(&balace); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) Withdraw(ctx context.Context, amount float64, id int) error { //transaction
	tx, err := p.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	query := `UPDATE public."Account" SET balance = balance - $1 WHERE id = $2  RETURNING balance`

	row := tx.QueryRow(ctx, query, int(amount*100), id)

	var balace int

	if err := row.Scan(&balace); err != nil {
		return err
	}

	return nil
}
