package kafka

import (
	"context"
	"encoding/json"
	"github.com/novabankapp/common.infrastructure/tracing"
	"github.com/novabankapp/wallet.api/infrastructure/kafka/messages"
	walletCommands "github.com/novabankapp/wallet.application/commands"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
)

func (s *readerMessageProcessor) processUserAccountCreated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "readerMessageProcessor.processProductDeleted")
	defer span.Finish()

	msg := messages.UserAccountCreatedMessage{}
	err := json.Unmarshal(m.Value, msg)
	if err != nil {
		s.commitErrMessage(ctx, r, m)
		return
	}
	id := uuid.NewV4().String()
	walletId := uuid.NewV4().String()
	amount := decimal.NewFromFloat(0.00)
	command := walletCommands.NewCreateWalletCommand(id, amount, "", msg.UserId, msg.AccountId, walletId)
	err = s.walletService.Commands.CreateWallet.Handle(ctx, command)
	if err != nil {
		s.commitErrMessage(ctx, r, m)
		return
	}
}
