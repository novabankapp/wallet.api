package resources

type UnlockWalletRequest struct {
	WalletId    string `json:"wallet_id"`
	Description string `json:"description"`
}
