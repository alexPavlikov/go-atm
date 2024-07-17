package locations

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	models "github.com/alexPavlikov/go-atm/internal/domain"
	"github.com/alexPavlikov/go-atm/internal/server/service"
)

type Handler struct {
	service service.Services
}

func New(service service.Services) *Handler {
	return &Handler{
		service: service,
	}
}

type emptyRequest struct{}

type emptyResponse struct{}

type accountsHandlerResponse struct {
	ID int
}

func (h *Handler) AccountsHandler(r *http.Request, data UserAccountHandlers) (resp accountsHandlerResponse, err error) {

	ctx := r.Context()

	var acc = UserAccountHandlers{
		ID:      data.ID,
		Name:    data.Name,
		Balance: data.Balance,
		Service: h.service,
		Ctx:     ctx,
	}

	if err := acc.Create(); err != nil {
		return accountsHandlerResponse{}, err
	}

	resp.ID = acc.ID

	return resp, nil
}

func (h *Handler) DepositAccountsHandler(r *http.Request, data models.DepositRequest) (emptyResponse, error) {
	ctx := r.Context()
	id := ctx.Value("Account_id")

	var acc = UserAccountHandlers{
		ID:      id.(int),
		Service: h.service,
		Ctx:     ctx,
	}

	if err := acc.Deposite(data.Deposit); err != nil {
		return emptyResponse{}, err
	}

	return emptyResponse{}, nil
}

func (h *Handler) WithdrawAccountsHandler(r *http.Request, data models.WithdrawRequest) (emptyResponse, error) {
	ctx := r.Context()
	id := ctx.Value("Account_id")

	var acc = UserAccountHandlers{
		ID:      id.(int),
		Service: h.service,
		Ctx:     ctx,
	}

	if err := acc.Withdraw(data.Withdraw); err != nil {
		return emptyResponse{}, err
	}

	return emptyResponse{}, nil
}

func (h *Handler) BalanceAccountsHandler(r *http.Request, data emptyRequest) (br models.BalanceResponse, err error) {
	ctx := r.Context()
	id := ctx.Value("Account_id")

	var acc = UserAccountHandlers{
		ID:      id.(int),
		Service: h.service,
		Ctx:     ctx,
	}

	br.Balance = acc.GetBalance()

	return br, nil
}

type UserAccountHandlers struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`

	Service service.Services
	Ctx     context.Context
}

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

type UserAccount interface {
	Create() error
}

func (acc *UserAccountHandlers) Create() error {

	var srvAcc = models.UserAccount{
		ID:      acc.ID,
		Name:    acc.Name,
		Balance: acc.Balance,
	}

	var err error
	acc.ID, err = acc.Service.AddAccount(acc.Ctx, srvAcc)
	if err != nil {
		return fmt.Errorf("create user error: %w", err)
	}
	return nil
}

func (acc *UserAccountHandlers) Deposite(amount float64) error {
	if err := acc.Service.UpDeposite(acc.Ctx, amount, acc.ID); err != nil {
		return fmt.Errorf("func deposite err: %w", err)
	}

	return nil
}

func (acc *UserAccountHandlers) Withdraw(amount float64) error {
	serviceAccount, err := acc.Service.GetAccount(acc.Ctx, acc.ID)
	if err != nil {
		return fmt.Errorf("func withdraw get account err: %w", err)
	}

	if int(serviceAccount.Balance*100) < int(amount*100) {
		return errors.New("not enough money on the balance")
	}

	if err := acc.Service.Withdraw(acc.Ctx, amount, acc.ID); err != nil {
		return fmt.Errorf("func withdraw err: %w", err)
	}

	return nil
}

func (acc *UserAccountHandlers) GetBalance() float64 {
	serviceAccount, err := acc.Service.GetAccount(acc.Ctx, acc.ID)
	if err != nil {
		return -1
	}
	return serviceAccount.Balance
}
