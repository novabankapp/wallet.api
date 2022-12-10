package kafka

import (
	"context"
	"encoding/json"
	"github.com/novabankapp/common.infrastructure/tracing"
	walletResources "github.com/novabankapp/wallet.api/functions/wallets/resources"
	"github.com/novabankapp/wallet.api/infrastructure/kafka/messages"
	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
)

func (s *readerMessageProcessor) processUserAccountCreated(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "readerMessageProcessor.processUserAccountCreated")
	defer span.Finish()

	msg := messages.UserAccountCreatedMessage{}
	err := json.Unmarshal(m.Value, msg)
	if err != nil {
		s.commitErrMessage(ctx, r, m)
		return
	}

	amount := decimal.NewFromFloat(0.00)
	_, err = s.walletService.CreateWallet(ctx, walletResources.CreateWalletRequest{
		AccountId:   msg.AccountId,
		UserId:      msg.UserId,
		Description: "",
		Amount:      &amount,
	})
	if err != nil {
		s.commitErrMessage(ctx, r, m)
		return
	}
}
