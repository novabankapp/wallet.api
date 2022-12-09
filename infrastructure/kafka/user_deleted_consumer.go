package kafka

import (
	"context"
	"encoding/json"
	"github.com/novabankapp/common.infrastructure/tracing"
	"github.com/novabankapp/wallet.api/infrastructure/kafka/messages"
	walletCommands "github.com/novabankapp/wallet.application/commands"
	walletQueries "github.com/novabankapp/wallet.application/queries"
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
	results, _, err := s.walletService.Queries.GetUserWalletsByID.Handle(ctx, &walletQueries.GetUserWalletsByIDQuery{
		UserID: msg.UserId,
	}, 1, []byte(""))
	if err != nil {
		s.commitErrMessage(ctx, r, m)
		return
	}
	for _, v := range *results {
		command := walletCommands.NewDeleteWalletCommand(v.WalletID, v.WalletID, "")
		s.walletService.Commands.DeleteWalletCommand.Handle(ctx, command)
	}
}
