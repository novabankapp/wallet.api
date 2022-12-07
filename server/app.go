package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/wallet.api/config"
	readerKafka "github.com/novabankapp/wallet.api/infrastructure/kafka"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type App interface {
	Run(ctx context.Context) error
}

type app struct {
	httpServer *http.Server
	log        logger.Logger
	cfg        *config.Config
}

func NewApp(server *http.Server, log logger.Logger, config *config.Config) App {
	return &app{
		server,
		log,
		config,
	}
}

var kafkaConn *kafka.Conn

// Run @title Novabank Wallet API
// @version 1.0
// @description This is a Wallet API.
// @BasePath /
func (a *app) Run(ctx context.Context) error {

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	readerMessageProcessor := readerKafka.NewReaderMessageProcessor(a.log, a.cfg)
	a.log.Info("Starting Reader Kafka consumers")
	cg := kafkaClient.NewConsumerGroup(a.cfg.Kafka.Brokers, a.cfg.Kafka.GroupID, a.log)
	go cg.ConsumeTopic(ctx, a.getConsumerGroupTopics(), readerKafka.PoolSize, readerMessageProcessor.ProcessMessages)

	if err := a.connectKafkaBrokers(ctx); err != nil {
		fmt.Println(err.Error())
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer kafkaConn.Close()

	if a.cfg.Kafka.InitTopics {
		a.initKafkaTopics(ctx)
	}

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
