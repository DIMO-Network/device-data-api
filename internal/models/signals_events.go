package models

type SignalsEvents struct {
	Name       string `boil:"name" json:"name" toml:"name" yaml:"name"`
	TotalCount int64
}

type DateIDItem struct {
	DateID        string `boil:"date_id" json:"date_id" toml:"date_id" yaml:"date_id"`
	IntegrationID string `boil:"integration_id" json:"integration_id" toml:"integration_id" yaml:"integration_id"`
}

type SignalsEventsUserDevices struct {
	IntegrationID              string `boil:"integration_id" json:"integration_id" toml:"integration_id" yaml:"integration_id"`
	PowerTrainType             string `boil:"power_train_type" json:"power_train_type" toml:"power_train_type" yaml:"power_train_type"`
	TotalCount                 int64
	TotalDeviceDefinitionCount int64
}
