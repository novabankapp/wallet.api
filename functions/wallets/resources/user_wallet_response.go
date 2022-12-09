package resources

type UserWalletsResponse struct {
	Cursor   *string
	PageSize *int
	Wallets  []WalletResponse
}
