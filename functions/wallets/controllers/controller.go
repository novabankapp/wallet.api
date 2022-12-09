package controllers

import (
	"context"
	"github.com/novabankapp/common.application/utilities/cryptography"
	"github.com/novabankapp/wallet.api/functions/common/resources"
	walletResources "github.com/novabankapp/wallet.api/functions/wallets/resources"
	walletQueries "github.com/novabankapp/wallet.application/queries"
	walletServices "github.com/novabankapp/wallet.application/services"
)

type WalletController interface {
}

type walletController struct {
	service      walletServices.WalletService
	cryptography cryptography.Cryptography
}

func newWalletController(service walletServices.WalletService, cryptography cryptography.Cryptography) WalletController {
	return &walletController{service: service, cryptography: cryptography}
}

func (w *walletController) GetWalletById(ctx context.Context, walletId string) (*walletResources.WalletResponse, error) {

	walletDto, err := w.service.Queries.GetWalletByID.Handle(ctx, &walletQueries.GetWalletByIDQuery{
		ID: walletId,
	})
	if err != nil {
		return nil, err
	}
	return &walletResources.WalletResponse{
		ID:               walletDto.ID,
		WalletID:         walletDto.WalletID,
		UserId:           walletDto.Wallet.UserId,
		AccountId:        walletDto.Wallet.AccountId,
		Balance:          walletDto.Wallet.Balance,
		AvailableBalance: walletDto.Wallet.AvailableBalance,
		IsLocked:         walletDto.WalletState.IsLocked,
		IsBlacklisted:    walletDto.WalletState.IsBlacklisted,
		IsDeleted:        walletDto.WalletState.IsDeleted,
		CreatedAt:        walletDto.Wallet.CreatedAt,
	}, nil
}
func (w *walletController) GetWalletsByUserId(ctx context.Context, userId string, pageData resources.PaginationData) (*walletResources.UserWalletsResponse, error) {
	var page []byte
	if pageData.PageCursor != nil {
		var err error
		page, err = w.cryptography.DecryptString(*pageData.PageCursor, nil)
		if err != nil {
			return nil, err
		}
	}
	pageSize := pageData.PageSize
	if pageData.PageSize == nil {
		*pageSize = 1
	}
	results, pageState, err := w.service.Queries.GetUserWalletsByID.Handle(ctx,
		&walletQueries.GetUserWalletsByIDQuery{
			UserID: userId,
		}, *pageSize, page)
	cursor, err := w.cryptography.EncryptAsString(pageState, nil)
	if err != nil {
		return nil, err
	}
	var wallets = make([]walletResources.WalletResponse, len(*results))
	for _, w := range *results {
		wallets = append(wallets,
			walletResources.WalletResponse{
				ID:               w.ID,
				WalletID:         w.WalletID,
				UserId:           w.Wallet.UserId,
				AccountId:        w.Wallet.AccountId,
				Balance:          w.Wallet.Balance,
				AvailableBalance: w.Wallet.AvailableBalance,
				IsLocked:         w.WalletState.IsLocked,
				IsBlacklisted:    w.WalletState.IsBlacklisted,
				IsDeleted:        w.WalletState.IsDeleted,
				CreatedAt:        w.Wallet.CreatedAt,
			})
	}
	return &walletResources.UserWalletsResponse{
		Cursor:   &cursor,
		PageSize: pageSize,
		Wallets:  wallets,
	}, nil
}
