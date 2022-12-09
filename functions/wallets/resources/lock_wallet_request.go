package resources

type LockWalletRequest struct {
	WalletId    string `json:"wallet_id"`
	Description string `json:"description"`
}
