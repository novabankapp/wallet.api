package resources

import "github.com/shopspring/decimal"

type CreateWalletRequest struct {
	AccountId   string           `json:"account_id"`
	UserId      string           `json:"user_id"`
	Description string           `json:"description"`
	Amount      *decimal.Decimal `json:"amount"`
}
