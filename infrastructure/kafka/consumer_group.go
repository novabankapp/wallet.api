package kafka

import (
	"context"
	"sync"

	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/wallet.api/config"
	"github.com/segmentio/kafka-go"
)

const (
	PoolSize = 30
)

type readerMessageProcessor struct {
	log logger.Logger
	cfg *config.Config
}

func NewReaderMessageProcessor(log logger.Logger,
	cfg *config.Config) *readerMessageProcessor {
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

		}
	}
}

func (p *readerMessageProcessor) logProcessMessage(m kafka.Message, workerID int) {
	p.log.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time)
}
