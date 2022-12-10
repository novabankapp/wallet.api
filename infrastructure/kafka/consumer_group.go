package kafka

import (
	"context"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/wallet.api/config"
	walletServices "github.com/novabankapp/wallet.application/services"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

type readerMessageProcessor struct {
	log           logger.Logger
	cfg           *config.Config
	walletService walletServices.WalletService
}

func NewReaderMessageProcessor(log logger.Logger,
	cfg *config.Config, service walletServices.WalletService) *readerMessageProcessor {
	return &readerMessageProcessor{log: log, cfg: cfg}
}

func (p *readerMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			p.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		p.logProcessMessage(m, workerID)

		switch m.Topic {
		case p.cfg.Kafka.KafkaTopics.UserCreated.TopicName:
			p.processUserCreated(ctx, r, m)
		case p.cfg.Kafka.KafkaTopics.AccountCreated.TopicName:
			p.processUserAccountCreated(ctx, r, m)
		case p.cfg.Kafka.KafkaTopics.UserDeleted.TopicName:
			p.processUserDeleted(ctx, r, m)

		}
	}
}
