package resources

type UserWalletsResponse struct {
	Cursor   *string          `json:"cursor"`
	PageSize *int             `json:"page_size"`
	Wallets  []WalletResponse `json:"wallets"`
}
