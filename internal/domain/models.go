package models

type UserAccount struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type DepositRequest struct {
	Deposit float64 `json:"deposit"`
}

type WithdrawRequest struct {
	Withdraw float64 `json:"withdraw"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}
