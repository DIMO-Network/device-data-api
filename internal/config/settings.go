package config

import (
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
)

// Settings contains the application config
type Settings struct {
	Environment                    string `yaml:"ENVIRONMENT"`
	Port                           string `yaml:"PORT"`
	GRPCPort                       string `yaml:"GRPC_PORT"`
	LogLevel                       string `yaml:"LOG_LEVEL"`
	ServiceName                    string `yaml:"SERVICE_NAME"`
	JwtKeySetURL                   string `yaml:"JWT_KEY_SET_URL"`
	DeploymentBaseURL              string `yaml:"DEPLOYMENT_BASE_URL"`
	ElasticSearchAnalyticsHost     string `yaml:"ELASTIC_SEARCH_ANALYTICS_HOST"`
	ElasticSearchAnalyticsUsername string `yaml:"ELASTIC_SEARCH_ANALYTICS_USERNAME"`
	ElasticSearchAnalyticsPassword string `yaml:"ELASTIC_SEARCH_ANALYTICS_PASSWORD"`
	DeviceDataIndexName            string `yaml:"DEVICE_DATA_INDEX_NAME"`
	DevicesAPIGRPCAddr             string `yaml:"DEVICES_APIGRPC_ADDR"`
	EmailHost                      string `yaml:"EMAIL_HOST"`
	EmailPort                      string `yaml:"EMAIL_PORT"`
	EmailUsername                  string `yaml:"EMAIL_USERNAME"`
	EmailPassword                  string `yaml:"EMAIL_PASSWORD"`
	EmailFrom                      string `yaml:"EMAIL_FROM"`
	AWSEndpoint                    string `yaml:"AWS_ENDPOINT"`
	AWSBucketName                  string `yaml:"AWS_BUCKET_NAME"`
	AWSAccessKeyID                 string `yaml:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey             string `yaml:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion                      string `yaml:"AWS_REGION"`
	UsersAPIGRPCAddr               string `yaml:"USERS_API_GRPC_ADDR"`
	NATSURL                        string `yaml:"NATS_URL"`
	NATSStreamName                 string `yaml:"NATS_STREAM_NAME"`
	NATSDataDownloadSubject        string `yaml:"NATS_DATA_DOWNLOAD_SUBJECT"`
	NATSAckTimeout                 string `yaml:"NATS_ACK_TIMEOUT"`
	NATSDurableConsumer            string `yaml:"NATS_DURABLE_CONSUMER"`
	MaxFileSize                    int    `yaml:"MAX_AWS_FILE_SIZE"`

	EnablePrivileges          bool        `yaml:"ENABLE_PRIVILEGES"`
	TokenExchangeJWTKeySetURL string      `yaml:"TOKEN_EXCHANGE_JWK_KEY_SET_URL"`
	VehicleNFTAddress         string      `yaml:"VEHICLE_NFT_ADDRESS"`
	DeviceDefinitionsGRPCAddr string      `yaml:"DEVICE_DEFINITIONS_GRPC_ADDR"`
	DB                        db.Settings `yaml:"DB"`
	AutoPiAPIToken            string      `yaml:"AUTO_PI_API_TOKEN"`
	AutoPiAPIURL              string      `yaml:"AUTO_PI_API_URL"`
	DeviceStatusTopic         string      `yaml:"DEVICE_STATUS_TOPIC"`
	KafkaBrokers              string      `yaml:"KAFKA_BROKERS"`
	EventsTopic               string      `yaml:"EVENTS_TOPIC"`
}

func (s *Settings) IsKafkaEnabled(logger *zerolog.Logger) bool {
	if s.KafkaBrokers == "" {
		logger.Info().Msg("KAFKA_BROKERS is not set, any dependencies should be disabled")
		return false
	}
	return true
}

func (s *Settings) IsWebAPIEnabled(logger *zerolog.Logger) bool {
	if s.Port == "" {
		logger.Info().Msg("PORT is not set, web api should be disabled")
		return false
	}
	return true
}
