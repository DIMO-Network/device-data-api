package services

// API wrapper to call device-data-api to get the userDevices associated with a userId and other use cases.
// ideally this would then be with grpc

//go:generate mockgen -source device_api_service.go -destination mocks/device_api_service_mock.go
type DeviceAPIService interface {
	GetUserDevices(userID string) ([]UserDeviceFull, error)
	UserDeviceBelongsToUserId(userID, userDeviceID string) (bool, error)
}

func NewDeviceAPIService(deviceApiUrl string) DeviceAPIService {
	return &deviceAPIService{deviceApiUrl: deviceApiUrl}
}

type deviceAPIService struct {
	deviceApiUrl string
}

func (das deviceAPIService) GetUserDevices(userID string) ([]UserDeviceFull, error) {
	return nil, nil
}

func (das deviceAPIService) UserDeviceBelongsToUserId(userID, userDeviceID string) (bool, error) {
	return true, nil
}

// UserDeviceFull represents object user's see on frontend for listing of their devices
type UserDeviceFull struct {
	ID             string  `json:"id"`
	VIN            *string `json:"vin"`
	VINConfirmed   bool    `json:"vinConfirmed"`
	Name           *string `json:"name"`
	CustomImageURL *string `json:"customImageUrl"`
	CountryCode    *string `json:"countryCode"`
}
