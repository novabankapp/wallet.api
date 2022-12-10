package services

import (
	"context"
	"github.com/novabankapp/payment.application/dtos"
	paymentServices "github.com/novabankapp/payment.application/services"
	"github.com/shopspring/decimal"
)

type MoneyTransferService interface {
}

type moneyTransferService struct {
	transferService paymentServices.MoneyTransferService
}

func NewMoneyTransferService(transferService paymentServices.MoneyTransferService) MoneyTransferService {
	return &moneyTransferService{
		transferService: transferService,
	}
}

func (t *moneyTransferService) TransferFromWalletToWallet(
	ctx context.Context,
	creditWalletAggregateId string,
	amount decimal.Decimal,
	debitWalletAggregateId,
	description string,
) (result bool, error error) {

	return t.transferService.TransferFromWalletToWallet(ctx, creditWalletAggregateId, amount, debitWalletAggregateId, description)

}

func (t *moneyTransferService) ReceiveFromServiceToWallet(
	ctx context.Context,
	fromNovaSuspenseWalletAggregateId,
	toWalletAggregateId,
	description string,
	amount decimal.Decimal,
) (result *string, error error) {
	return t.transferService.ReceiveFromServiceToWallet(ctx, fromNovaSuspenseWalletAggregateId, toWalletAggregateId, description, amount)
}

func (t *moneyTransferService) TransferToServiceFromWallet(
	ctx context.Context,
	from dtos.AccountDto,
	to dtos.AccountDto,
	amount decimal.Decimal,
	toNovaSuspenseWalletAggregateId,
	description string,
) (result *string, error error) {
	return t.transferService.TransferToServiceFromWallet(ctx, from, to, amount, toNovaSuspenseWalletAggregateId, description)
}
