package controllers

import (
	"context"
	"github.com/novabankapp/wallet.api/functions/common/resources"
	walletResources "github.com/novabankapp/wallet.api/functions/wallets/resources"

	walletServices "github.com/novabankapp/wallet.api/functions/wallets/services"
)

type WalletController interface {
	GetWalletById(ctx context.Context, walletId string) (*walletResources.WalletResponse, error)
	CreateWallet(ctx context.Context, req walletResources.CreateWalletRequest) (bool, error)
	BlockWallet(ctx context.Context, req walletResources.BlockWalletRequest) (bool, error)
	UnblockWallet(ctx context.Context, req walletResources.UnblockWalletRequest) (bool, error)
	DeleteWallet(ctx context.Context, req walletResources.DeleteWalletRequest) (bool, error)
	LockWallet(ctx context.Context, req walletResources.LockWalletRequest) (bool, error)
	UnlockWallet(ctx context.Context, req walletResources.UnlockWalletRequest) (bool, error)
	DebitWallet(ctx context.Context, req walletResources.DebitAccountRequest) (bool, error)
	CreditWallet(ctx context.Context, req walletResources.CreditAccountRequest) (bool, error)
	GetWalletsByUserId(ctx context.Context, userId string, pageData resources.PaginationData) (*walletResources.UserWalletsResponse, error)
}

type walletController struct {
	service walletServices.WalletService
}

func NewWalletController(service walletServices.WalletService) WalletController {
	return &walletController{service: service}
}

func (w *walletController) GetWalletById(ctx context.Context, walletId string) (*walletResources.WalletResponse, error) {

	return w.service.GetWalletById(ctx, walletId)
}
func (w *walletController) GetWalletsByUserId(ctx context.Context, userId string, pageData resources.PaginationData) (*walletResources.UserWalletsResponse, error) {
	return w.service.GetWalletsByUserId(ctx, userId, pageData)
}
func (w *walletController) CreateWallet(ctx context.Context, req walletResources.CreateWalletRequest) (bool, error) {
	return w.service.CreateWallet(ctx, req)
}
func (w *walletController) BlockWallet(ctx context.Context, req walletResources.BlockWalletRequest) (bool, error) {
	return w.service.BlockWallet(ctx, req)
}
func (w *walletController) UnblockWallet(ctx context.Context, req walletResources.UnblockWalletRequest) (bool, error) {
	return w.service.UnblockWallet(ctx, req)
}
func (w *walletController) DeleteWallet(ctx context.Context, req walletResources.DeleteWalletRequest) (bool, error) {
	return w.service.DeleteWallet(ctx, req)
}

func (w *walletController) LockWallet(ctx context.Context, req walletResources.LockWalletRequest) (bool, error) {
	return w.service.LockWallet(ctx, req)
}
func (w *walletController) UnlockWallet(ctx context.Context, req walletResources.UnlockWalletRequest) (bool, error) {
	return w.service.UnlockWallet(ctx, req)
}

func (w *walletController) DebitWallet(ctx context.Context, req walletResources.DebitAccountRequest) (bool, error) {
	return w.service.DebitWallet(ctx, req)
}
func (w *walletController) CreditWallet(ctx context.Context, req walletResources.CreditAccountRequest) (bool, error) {
	return w.service.CreditWallet(ctx, req)

}
