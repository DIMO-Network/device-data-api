package services

import (
	"context"
	"fmt"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"google.golang.org/grpc"
)

//go:generate mockgen -source device_api_service.go -destination mocks/device_api_service_mock.go
type DeviceAPIService interface {
	ListUserDevicesForUser(ctx context.Context, userID string) (*pb.ListUserDevicesForUserResponse, error)
	GetUserDevice(ctx context.Context, userDeviceID string) (*pb.UserDevice, error)
	UserDeviceBelongsToUserID(ctx context.Context, userID, userDeviceID string) (bool, error)
	GetUserDeviceByTokenID(ctx context.Context, tokenID int64) (*pb.UserDevice, error)
}

// NewDeviceAPIService API wrapper to call device-data-api to get the userDevices associated with a userId over grpc
func NewDeviceAPIService(devicesConn *grpc.ClientConn) DeviceAPIService {
	return &deviceAPIService{devicesConn: devicesConn}
}

type deviceAPIService struct {
	devicesConn *grpc.ClientConn
}

func (das *deviceAPIService) ListUserDevicesForUser(ctx context.Context, userID string) (*pb.ListUserDevicesForUserResponse, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("user id was empty - invalid")
	}
	deviceClient := pb.NewUserDeviceServiceClient(das.devicesConn)

	devicesForUser, err := deviceClient.ListUserDevicesForUser(ctx, &pb.ListUserDevicesForUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	return devicesForUser, nil
}

func (das *deviceAPIService) UserDeviceBelongsToUserID(ctx context.Context, userID, userDeviceID string) (bool, error) {
	device, err := das.GetUserDevice(ctx, userDeviceID)
	if err != nil {
		return false, err
	}
	return device.UserId == userID, nil
}

func (das *deviceAPIService) GetUserDevice(ctx context.Context, userDeviceID string) (*pb.UserDevice, error) {
	if len(userDeviceID) == 0 {
		return nil, fmt.Errorf("user device id was empty - invalid")
	}
	deviceClient := pb.NewUserDeviceServiceClient(das.devicesConn)

	userDevice, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{
		Id: userDeviceID,
	})
	if err != nil {
		return nil, err
	}

	return userDevice, nil
}

func (das *deviceAPIService) GetUserDeviceByTokenID(ctx context.Context, tokenID int64) (*pb.UserDevice, error) {
	deviceClient := pb.NewUserDeviceServiceClient(das.devicesConn)

	userDevice, err := deviceClient.GetUserDeviceByTokenId(ctx, &pb.GetUserDeviceByTokenIdRequest{
		TokenId: tokenID,
	})
	if err != nil {
		return nil, err
	}

	return userDevice, nil
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
