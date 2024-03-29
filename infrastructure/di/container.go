package di

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gocql/gocql"
	"github.com/novabankapp/common.application/services/message_queue"
	"github.com/novabankapp/common.application/utilities/cryptography"
	es "github.com/novabankapp/common.data/eventstore"
	store "github.com/novabankapp/common.data/eventstore/store"
	baseRepository "github.com/novabankapp/common.data/repositories/base/cassandra"
	"github.com/novabankapp/common.infrastructure/eventstoredb"
	"github.com/novabankapp/common.infrastructure/kafka"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/common.notifier/email"
	"github.com/novabankapp/common.notifier/sms"
	paymentServices "github.com/novabankapp/payment.application/services"
	integrationServices "github.com/novabankapp/payment.application/services/integrations"
	localConfig "github.com/novabankapp/wallet.api/config"
	allControllers "github.com/novabankapp/wallet.api/controllers"
	paymentControllers "github.com/novabankapp/wallet.api/functions/payments/controllers"
	paymentServicesLocal "github.com/novabankapp/wallet.api/functions/payments/services"
	walletControllers "github.com/novabankapp/wallet.api/functions/wallets/controllers"
	walletServicesLocal "github.com/novabankapp/wallet.api/functions/wallets/services"
	"github.com/novabankapp/wallet.api/middlewares"
	"github.com/novabankapp/wallet.api/server"
	walletServices "github.com/novabankapp/wallet.application/services"
	noSqlEntities "github.com/novabankapp/wallet.data/es/models"
	"github.com/novabankapp/wallet.data/migrations"
	"github.com/scylladb/gocqlx/v2"
	"go.uber.org/dig"
	"net/http"
)

//var container *dig.Container = dig.New()

func BuildContainer() *dig.Container {
	container := dig.New()
	err := container.Provide(func() (*localConfig.Config, error) {

		return localConfig.InitConfig()
	})
	if err != nil {
		return nil
	}
	loggerError := container.Provide(func(config *localConfig.Config) logger.Logger {

		appLogger := logger.NewAppLogger(config.Logger)
		appLogger.InitLogger()
		appLogger.WithName("UserManagementAPI")
		return appLogger
	})
	if loggerError != nil {
		return nil
	}
	eventStoreDbError := container.Provide(func(cfg eventstoredb.EventStoreConfig) (*esdb.Client, error) {
		settings, err := esdb.ParseConnectionString(cfg.ConnectionString)
		if err != nil {
			return nil, err
		}
		return esdb.NewClient(settings)

	})
	if eventStoreDbError != nil {
		return nil
	}
	container.Provide(func(log logger.Logger, db *esdb.Client) es.AggregateStore {
		return store.NewAggregateStore(log, db)
	})
	container.Provide(func(log logger.Logger, db *esdb.Client) es.EventStore {
		return store.NewEventStore(log, db)
	})
	cassandraSessionError := container.Provide(func(config *localConfig.Config) (session *gocqlx.Session, err error) {

		cluster := gocql.NewCluster(config.Cassandra.Addresses...)
		cluster.ProtoVersion = config.Cassandra.ProtoVersion
		cluster.Keyspace = config.Cassandra.Keyspace
		//cluster.ConnectTimeout = config.Cassandra.Timeout
		//cluster.Authenticator = gocql.PasswordAuthenticator{Username: config.Cassandra.Username, Password: config.Cassandra.Password}
		ss, ee := cluster.CreateSession()

		s, e := gocqlx.WrapSession(ss, ee)
		if e != nil {
			return nil, e
		}
		migrations.InitCassandra(&s)
		return &s, e

	})
	if cassandraSessionError != nil {
		return nil
	}

	kafkaProducerErr := container.Provide(func(appLogger logger.Logger, config *localConfig.Config) kafka.Producer {

		return kafka.NewProducer(appLogger, config.Kafka.Brokers)
	})
	if kafkaProducerErr != nil {
		return nil
	}
	kafkaTopicsError := container.Provide(func(config *localConfig.Config) *kafkaClient.KafkaTopics {

		return &kafkaClient.KafkaTopics{
			UserCreated:         config.Kafka.KafkaTopics.UserCreated,
			UserUpdated:         config.Kafka.KafkaTopics.UserUpdated,
			UserDeleted:         config.Kafka.KafkaTopics.UserDeleted,
			ContactDeleted:      config.Kafka.KafkaTopics.ContactDeleted,
			ContactUpdated:      config.Kafka.KafkaTopics.ContactUpdated,
			UserPasswordChanged: config.Kafka.KafkaTopics.UserPasswordChanged,
			UserLocked:          config.Kafka.KafkaTopics.UserLocked,
			UserLoggedIn:        config.Kafka.KafkaTopics.UserLoggedIn,
			AccountCreated:      config.Kafka.KafkaTopics.AccountCreated,
			AccountLocked:       config.Kafka.KafkaTopics.AccountLocked,
			AccountDeactivated:  config.Kafka.KafkaTopics.AccountDeactivated,
			AccountUnlocked:     config.Kafka.KafkaTopics.AccountUnlocked,
			AccountActivated:    config.Kafka.KafkaTopics.AccountActivated,
		}
	})
	if kafkaTopicsError != nil {
		return nil
	}

	messageQueueErr := container.Provide(func(producer kafka.Producer) message_queue.MessageQueue {
		return message_queue.NewKafkaMessageQueue(producer)
	})
	if messageQueueErr != nil {
		return nil
	}

	container.Provide(func(config *localConfig.Config) (sms.SMSService, error) {
		//return sms.NewSMPPService(config.SMPP)
		return sms.NewMockSMSService(), nil
	})
	container.Provide(func(config *localConfig.Config) email.MailService {
		//return email.NewSmtpService(config.SMTP)
		return email.NewMockMailService()
	})

	container.Provide(func(
		logger logger.Logger,
		config *localConfig.Config,
		session *gocqlx.Session,
		aggregateStore es.AggregateStore,
		topics *kafka.KafkaTopics,
		messageQueue message_queue.MessageQueue,
	) *walletServices.WalletService {
		walletProjectionRepo := baseRepository.NewCassandraRepository[noSqlEntities.WalletProjection](session, "", config.Cassandra.Timeout)
		return walletServices.NewWalletService(logger, aggregateStore, walletProjectionRepo, topics, messageQueue)
	})
	container.Provide(func() cryptography.Cryptography {
		return cryptography.New("novabank")
	})
	container.Provide(func(walletService walletServices.WalletService, cryptography cryptography.Cryptography) walletServicesLocal.WalletService {
		return walletServicesLocal.NewWalletService(walletService, cryptography)
	})

	container.Provide(func() integrationServices.MoneyService {
		return integrationServices.NewMoneyService()
	})

	container.Provide(func(walletService walletServices.WalletService, moneyService integrationServices.MoneyService) paymentServices.MoneyTransferService {
		return paymentServices.NewMoneyTransferService(walletService, moneyService)
	})

	container.Provide(func(moneyTransferService paymentServices.MoneyTransferService) paymentServicesLocal.MoneyTransferService {
		return paymentServicesLocal.NewMoneyTransferService(moneyTransferService)
	})
	//controllers
	container.Provide(func(walletService walletServicesLocal.WalletService) walletControllers.WalletController {
		return walletControllers.NewWalletController(walletService)
	})

	container.Provide(func(moneyTransferService paymentServicesLocal.MoneyTransferService) paymentControllers.PaymentController {
		return paymentControllers.NewPaymentController(moneyTransferService)
	})

	newControllerErr := container.Provide(allControllers.NewControllers)
	if newControllerErr != nil {
		return nil
	}

	middlewareErr := container.Provide(middlewares.NewMiddlewares)
	if middlewareErr != nil {
		return nil
	}

	serverErr := container.Provide(func(config *localConfig.Config, controllers *allControllers.Controllers, middlewares middlewares.Middlewares, logger logger.Logger) *http.Server {
		return server.NewServer(config.Api.Address, config.Api.Port, *controllers, middlewares, logger)
	})
	if serverErr != nil {
		return nil
	}
	appErr := container.Provide(func(httpServer *http.Server, log logger.Logger, config *localConfig.Config, walletService walletServicesLocal.WalletService) server.App {
		return server.NewApp(httpServer, log, config, walletService)
	})
	if appErr != nil {
		return nil
	}

	return container
}
