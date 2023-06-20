package constants

const (
	AutoPiVendor   = "AutoPi"
	SmartCarVendor = "SmartCar"
	TeslaVendor    = "Tesla"
)

const (
	IntegrationTypeHardware string = "Hardware"
	IntegrationTypeAPI      string = "API"
)

const (
	IntegrationStyleAddon   string = "Addon"
	IntegrationStyleOEM     string = "OEM"
	IntegrationStyleWebhook string = "Webhook"
)

type RegionEnum string

const (
	AmericasRegion RegionEnum = "Americas"
	EuropeRegion   RegionEnum = "Europe"
)

func (r RegionEnum) String() string {
	return string(r)
}

// Enum values for UserDeviceAPIIntegrationStatus
const (
	UserDeviceAPIIntegrationStatusPending               string = "Pending"
	UserDeviceAPIIntegrationStatusPendingFirstData      string = "PendingFirstData"
	UserDeviceAPIIntegrationStatusActive                string = "Active"
	UserDeviceAPIIntegrationStatusFailed                string = "Failed"
	UserDeviceAPIIntegrationStatusDuplicateIntegration  string = "DuplicateIntegration"
	UserDeviceAPIIntegrationStatusAuthenticationFailure string = "AuthenticationFailure"
)
