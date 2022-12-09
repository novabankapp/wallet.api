package controllers

import (
	"context"
	"github.com/novabankapp/common.application/utilities/cryptography"
	"github.com/novabankapp/wallet.api/functions/common/resources"
	walletResources "github.com/novabankapp/wallet.api/functions/wallets/resources"
	walletCommands "github.com/novabankapp/wallet.application/commands"
	walletQueries "github.com/novabankapp/wallet.application/queries"
	walletServices "github.com/novabankapp/wallet.application/services"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
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
func (w *walletController) CreateWallet(ctx context.Context, req walletResources.CreateWalletRequest) {

	id := uuid.NewV4().String()
	walletId := uuid.NewV4().String()
	var amount decimal.Decimal
	if req.Amount != nil {
		amount = *req.Amount
	} else {
		amount = decimal.NewFromFloat(0.00)
	}
	command := walletCommands.NewCreateWalletCommand(id,
		amount, req.Description, req.UserId, req.AccountId, walletId)
	err := w.service.Commands.CreateWallet.Handle(ctx, command)
	if err != nil {
		return
	}
}
func (w *walletController) BlockWallet(ctx context.Context, req walletResources.BlockWalletRequest) (bool, error) {
	id := uuid.NewV4().String()
	command := walletCommands.NewBlockWalletCommand(id, req.WalletId, req.Description)
	err := w.service.Commands.BlockWalletCommand.Handle(ctx, command)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (w *walletController) UnblockWallet(ctx context.Context, req walletResources.UnblockWalletRequest) (bool, error) {
	id := uuid.NewV4().String()
	command := walletCommands.NewUnblockWalletCommand(id, req.WalletId, req.Description)
	err := w.service.Commands.UnblockWalletCommand.Handle(ctx, command)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (w *walletController) DeleteWallet(ctx context.Context, req walletResources.DeleteWalletRequest) (bool, error) {
	id := uuid.NewV4().String()
	command := walletCommands.NewDeleteWalletCommand(id, req.WalletId, req.Description)
	err := w.service.Commands.DeleteWalletCommand.Handle(ctx, command)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (w *walletController) LockWallet(ctx context.Context, req walletResources.LockWalletRequest) (bool, error) {
	id := uuid.NewV4().String()
	command := walletCommands.NewLockWalletCommand(id, req.WalletId, req.Description)
	err := w.service.Commands.LockWalletCommand.Handle(ctx, command)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (w *walletController) UnlockWallet(ctx context.Context, req walletResources.UnlockWalletRequest) (bool, error) {
	id := uuid.NewV4().String()
	command := walletCommands.NewUnlockWalletCommand(id, req.WalletId, req.Description)
	err := w.service.Commands.UnlockWalletCommand.Handle(ctx, command)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (w *walletController) DebitWallet(ctx context.Context, req walletResources.DebitAccountRequest) (bool, error) {

	command := walletCommands.NewDebitWalletCommand(req.DebitWalletID, req.CreditWalletID, req.Amount, req.Description)
	err := w.service.Commands.DebitWalletCommand.Handle(ctx, command)
	if err != nil {
		return false, err
	}
	return true, nil
}
func (w *walletController) CreditWallet(ctx context.Context, req walletResources.CreditAccountRequest) (bool, error) {

	command := walletCommands.NewCreditWalletCommand(req.CreditWalletID, req.DebitWalletID, req.Amount, req.Description)
	err := w.service.Commands.CreditWalletCommand.Handle(ctx, command)
	if err != nil {
		return false, err
	}
	return true, nil
}
