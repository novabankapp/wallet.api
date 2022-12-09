package resources

import "github.com/shopspring/decimal"

type CreditAccountRequest struct {
	CreditWalletID string          `json:"credit_wallet_id"`
	DebitWalletID  string          `json:"debit_wallet_id"`
	Description    string          `json:"description"`
	Amount         decimal.Decimal `json:"amount"`
}
