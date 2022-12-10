package controllers

import (
	paymentServices "github.com/novabankapp/wallet.api/functions/payments/services"
)

type PaymentController interface {
}

type paymentController struct {
	moneyTransferService paymentServices.MoneyTransferService
}

func NewPaymentController(moneyTransferService paymentServices.MoneyTransferService) PaymentController {
	return &paymentController{
		moneyTransferService: moneyTransferService,
	}
}
