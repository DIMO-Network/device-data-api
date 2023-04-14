package services

import (
	"context"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:generate mockgen -source device_api_service.go -destination mocks/device_api_service_mock.go
type DeviceAPIService interface {
	ListUserDevicesForUser(ctx context.Context, userID string) (*pb.ListUserDevicesForUserResponse, error)
	GetUserDevice(ctx context.Context, userDeviceID string) (*pb.UserDevice, error)
	UserDeviceBelongsToUserID(ctx context.Context, userID, userDeviceID string) (bool, error)
	GetUserDeviceByTokenID(ctx context.Context, tokenID int64) (*pb.UserDevice, error)
}

// NewDeviceAPIService API wrapper to call device-data-api to get the userDevices associated with a userId over grpc
func NewDeviceAPIService(devicesAPIGRPCAddr string) DeviceAPIService {
	return &deviceAPIService{devicesAPIGRPCAddr: devicesAPIGRPCAddr}
}

type deviceAPIService struct {
	devicesAPIGRPCAddr string
}

func (das *deviceAPIService) ListUserDevicesForUser(ctx context.Context, userID string) (*pb.ListUserDevicesForUserResponse, error) {
	conn, err := grpc.Dial(das.devicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	deviceClient := pb.NewUserDeviceServiceClient(conn)

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
	conn, err := grpc.Dial(das.devicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	deviceClient := pb.NewUserDeviceServiceClient(conn)

	userDevice, err := deviceClient.GetUserDevice(ctx, &pb.GetUserDeviceRequest{
		Id: userDeviceID,
	})
	if err != nil {
		return nil, err
	}

	return userDevice, nil
}

func (das *deviceAPIService) GetUserDeviceByTokenID(ctx context.Context, tokenID int64) (*pb.UserDevice, error) {
	conn, err := grpc.Dial(das.devicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	deviceClient := pb.NewUserDeviceServiceClient(conn)

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
