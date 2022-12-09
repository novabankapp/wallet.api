package resources

type UnblockWalletRequest struct {
	WalletId    string `json:"wallet_id"`
	Description string `json:"description"`
}
