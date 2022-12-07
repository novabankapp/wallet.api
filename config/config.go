package config

import (
	"flag"
	"fmt"
	"github.com/novabankapp/common.infrastructure/Cassandra"
	"github.com/novabankapp/common.infrastructure/constants"
	"github.com/novabankapp/common.infrastructure/eventstoredb"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/common.infrastructure/postgres"
	"github.com/novabankapp/common.infrastructure/redis"
	"github.com/novabankapp/common.notifier/email"
	"github.com/novabankapp/common.notifier/sms"
	localConstants "github.com/novabankapp/wallet.api/constants"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "Wallet API config path")
}

type Config struct {
	ServiceName     string                        `mapstructure:"serviceName" yaml:"serviceName"`
	Logger          *logger.Config                `mapstructure:"logger" yaml:"logger"`
	Postgresql      *postgres.Config              `mapstructure:"postgres" yaml:"postgres"`
	Kafka           *kafkaClient.Config           `mapstructure:"kafka" yaml:"kafka"`
	Redis           *redis.Config                 `mapstructure:"redis" yaml:"redis"`
	ServiceSettings ServiceSettings               `mapstructure:"serviceSettings" yaml:"serviceSettings"`
	Api             API                           `mapstructure:"api" yaml:"api"`
	Cassandra       Cassandra.Config              `mapstructure:"cassandra" yaml:"cassandra"`
	JwtToken        JwtToken                      `mapstructure:"jwtToken" yaml:"jwtToken"`
	SMPP            sms.SMPPConfig                `mapstructure:"smpp" yaml:"smpp"`
	SMTP            email.SMTPConfig              `mapstructure:"smtp" yaml:"smtp"`
	EventDBStore    eventstoredb.EventStoreConfig `mapstructure:"eventstore" yaml:"eventstore"`
	RequestTimeout  int
}
type JwtToken struct {
	SecretKey string `mapstructure:"secretKey" yaml:"secretkey"`
	Issuer    string `mapstructure:"issuer" yaml:"issuer"`
	ExpiresIn int    `mapstructure:"expiresInHours" yaml:"expiresInHours"`
}

type API struct {
	Port        string `mapstructure:"port" yaml:"port"`
	Address     string `mapstructure:"address" yaml:"address"`
	Development bool   `mapstructure:"development" yaml:"development"`
}

type ServiceSettings struct {
	RedisUserPrefixKey string `mapstructure:"redisUserPrefixKey" yaml:"redisUserPrefixKey"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		fmt.Println(configPathFromEnv)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config.yaml", getwd)
		}

	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}
	requestTimeout := os.Getenv(localConstants.REQUEST_TIMEOUT)
	if requestTimeout != "" {
		req, err := strconv.Atoi(requestTimeout)
		if err == nil {

			cfg.RequestTimeout = req
		} else {
			cfg.RequestTimeout = 5000
		}
	}

	apiPort := os.Getenv(localConstants.PORT)
	if apiPort != "" {
		cfg.Api.Port = apiPort
	}
	postgresHost := os.Getenv(constants.PostgresqlHost)
	if postgresHost != "" {
		cfg.Postgresql.Host = postgresHost
	}
	postgresPort := os.Getenv(constants.PostgresqlPort)
	if postgresPort != "" {
		cfg.Postgresql.Port = postgresPort
	}
	redisAddr := os.Getenv(constants.RedisAddr)
	if redisAddr != "" {
		cfg.Redis.Addr = redisAddr
	}

	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}

	return cfg, nil
}
