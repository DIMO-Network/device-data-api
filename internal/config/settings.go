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
}
