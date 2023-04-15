package services

import (
	"context"

	pb "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:generate mockgen -source device_definitions_service.go -destination mocks/device_definitions_service_mock.go
type DeviceDefinitionsAPIService interface {
	GetDeviceDefinition(ctx context.Context, id string) (*pb.GetDeviceDefinitionItemResponse, error)
}

func NewDeviceDefinitionsAPIService(deviceDefinitionsAPIGRPCAddr string) DeviceDefinitionsAPIService {
	return &deviceDefinitionsAPIService{deviceDefinitionsAPIGRPCAddr: deviceDefinitionsAPIGRPCAddr}
}

type deviceDefinitionsAPIService struct {
	deviceDefinitionsAPIGRPCAddr string
}

func (dda *deviceDefinitionsAPIService) GetDeviceDefinition(ctx context.Context, id string) (*pb.GetDeviceDefinitionItemResponse, error) {
	conn, err := grpc.Dial(dda.deviceDefinitionsAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	definitionsClient := pb.NewDeviceDefinitionServiceClient(conn)

	def, err := definitionsClient.GetDeviceDefinitionByID(ctx, &pb.GetDeviceDefinitionRequest{
		Ids: []string{id},
	})
	if err != nil {
		return nil, err
	}

	return def.DeviceDefinitions[0], nil
}

func (dda *deviceDefinitionsAPIService) GetDeviceStyle(ctx context.Context, id string) (*pb.DeviceStyle, error) {
	conn, err := grpc.Dial(dda.deviceDefinitionsAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	definitionsClient := pb.NewDeviceDefinitionServiceClient(conn)

	def, err := definitionsClient.GetDeviceStyleByID(ctx, &pb.GetDeviceStyleByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return def, nil
}
