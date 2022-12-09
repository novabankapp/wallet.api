package resources

type DeleteWalletRequest struct {
	WalletId    string `json:"wallet_id"`
	Description string `json:"description"`
}
