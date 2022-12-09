package resources

import (
	"github.com/shopspring/decimal"
	"time"
)

type WalletResponse struct {
	ID               string          `json:"id"`
	WalletID         string          `json:"wallet_id"`
	UserId           string          `json:"user_id"`
	AccountId        string          `json:"account_id"`
	Balance          decimal.Decimal `json:"balance"`
	AvailableBalance decimal.Decimal `json:"available_balance"`
	IsLocked         bool            `json:"is_locked"`
	IsBlacklisted    bool            `json:"is_blacklisted"`
	IsDeleted        bool            `json:"is_deleted"`
	CreatedAt        time.Time
}
