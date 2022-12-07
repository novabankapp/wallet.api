package server

import (
	"context"
	"fmt"
	"net"
	"strconv"

	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func (a *app) connectKafkaBrokers(ctx context.Context) (err error) {

	kafkaConn, err = kafkaClient.NewKafkaClient(ctx, *a.cfg.Kafka)
	if err != nil {

		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		fmt.Println(err.Error())
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	a.log.Infof("kafka connected to brokers: %+v", brokers)

	return err
}
func (a *app) getConsumerGroupTopics() []string {
	return []string{
		a.cfg.Kafka.KafkaTopics.UserLocked.TopicName,
	}
}
func (a *app) initKafkaTopics(ctx context.Context) {
	if kafkaConn == nil {
		fmt.Println("null here")
	}
	controller, err := kafkaConn.Controller()
	if err != nil {
		a.log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	a.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		a.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	a.log.Infof("established new kafka controller connection: %s", controllerURI)

	userLockedTopic := kafka.TopicConfig{
		Topic:             a.cfg.Kafka.KafkaTopics.UserLocked.TopicName,
		NumPartitions:     a.cfg.Kafka.KafkaTopics.UserLocked.Partitions,
		ReplicationFactor: a.cfg.Kafka.KafkaTopics.UserLocked.ReplicationFactor,
	}

	userDeletedTopic := kafka.TopicConfig{
		Topic:             a.cfg.Kafka.KafkaTopics.UserDeleted.TopicName,
		NumPartitions:     a.cfg.Kafka.KafkaTopics.UserDeleted.Partitions,
		ReplicationFactor: a.cfg.Kafka.KafkaTopics.UserDeleted.ReplicationFactor,
	}

	userUpdatedTopic := kafka.TopicConfig{
		Topic:             a.cfg.Kafka.KafkaTopics.UserUpdated.TopicName,
		NumPartitions:     a.cfg.Kafka.KafkaTopics.UserUpdated.Partitions,
		ReplicationFactor: a.cfg.Kafka.KafkaTopics.UserUpdated.ReplicationFactor,
	}

	userCreatedTopic := kafka.TopicConfig{
		Topic:             a.cfg.Kafka.KafkaTopics.UserCreated.TopicName,
		NumPartitions:     a.cfg.Kafka.KafkaTopics.UserCreated.Partitions,
		ReplicationFactor: a.cfg.Kafka.KafkaTopics.UserCreated.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		userLockedTopic,
		userDeletedTopic,
		userUpdatedTopic,
		userCreatedTopic,
	); err != nil {
		a.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	a.log.Infof("kafka topics created or already exists: %+v",
		[]kafka.TopicConfig{userLockedTopic,
			userDeletedTopic,
			userUpdatedTopic,
			userCreatedTopic})
}
