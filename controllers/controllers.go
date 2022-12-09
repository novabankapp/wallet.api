package controllers

import (
	walletControllers "github.com/novabankapp/wallet.api/functions/wallets/controllers"
)

type Controllers struct {
	walletController walletControllers.WalletController
}

func NewControllers(walletController walletControllers.WalletController) *Controllers {
	return &Controllers{
		walletController: walletController,
	}
}
