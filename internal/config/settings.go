package config

// Settings contains the application config
type Settings struct {
	Environment                    string `yaml:"ENVIRONMENT"`
	Port                           string `yaml:"PORT"`
	LogLevel                       string `yaml:"LOG_LEVEL"`
	ServiceName                    string `yaml:"SERVICE_NAME"`
	JwtKeySetURL                   string `yaml:"JWT_KEY_SET_URL"`
	DeploymentBaseURL              string `yaml:"DEPLOYMENT_BASE_URL"`
	ElasticSearchAnalyticsHost     string `yaml:"ELASTIC_SEARCH_ANALYTICS_HOST"`
	ElasticSearchAnalyticsUsername string `yaml:"ELASTIC_SEARCH_ANALYTICS_USERNAME"`
	ElasticSearchAnalyticsPassword string `yaml:"ELASTIC_SEARCH_ANALYTICS_PASSWORD"`
	DeviceDataIndexName            string `yaml:"DEVICE_DATA_INDEX_NAME"`
	DevicesAPIGRPCAddr             string `yaml:"DEVICES_APIGRPC_ADDR"`
	ElasticIndex                   string `yaml:"ELASTIC_INDEX"`
	EmailHost                      string `yaml:"EMAIL_HOST"`
	EmailPort                      string `yaml:"EMAIL_PORT"`
	EmailUsername                  string `yaml:"EMAIL_USERNAME"`
	EmailPassword                  string `yaml:"EMAIL_PASSWORD"`
	EmailFrom                      string `yaml:"EMAIL_FROM"`
	AWSBucketName                  string `yaml:"AWS_BUCKET_NAME"`
	AWSAccessKeyID                 string `yaml:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey             string `yaml:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion                      string `yaml:"AWS_REGION"`
	UsersAPIGRPCAddr               string `yaml:"USERS_API_GRPC_ADDR"`

	EnablePrivileges          bool   `yaml:"ENABLE_PRIVILEGES"`
	TokenExchangeJWTKeySetURL string `yaml:"TOKEN_EXCHANGE_JWK_KEY_SET_URL"`
	VehicleNFTAddress         string `yaml:"VEHICLE_NFT_ADDRESS"`
}
