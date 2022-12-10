package kafka

import (
	"context"
	"encoding/json"
	"github.com/novabankapp/common.infrastructure/tracing"
	"github.com/novabankapp/wallet.api/functions/common/resources"
	walletResources "github.com/novabankapp/wallet.api/functions/wallets/resources"
	"github.com/novabankapp/wallet.api/infrastructure/kafka/messages"
	"github.com/segmentio/kafka-go"
)

func (s *readerMessageProcessor) processUserDeleted(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "readerMessageProcessor.processUserDeleted")
	defer span.Finish()
	msg := messages.UserDeletedMessage{}
	err := json.Unmarshal(m.Value, msg)
	if err != nil {
		s.commitErrMessage(ctx, r, m)
		return
	}
	results, err := s.walletService.GetWalletsByUserId(ctx, msg.UserId, resources.PaginationData{
		PageSize:   nil,
		PageCursor: nil,
	})
	if err != nil {
		s.commitErrMessage(ctx, r, m)
		return
	}
	for _, v := range results.Wallets {
		s.walletService.DeleteWallet(ctx, walletResources.DeleteWalletRequest{
			WalletId: v.WalletID,
		})

	}
}
