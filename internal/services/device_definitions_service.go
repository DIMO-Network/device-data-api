package services

import (
	"context"
	"fmt"

	pb "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"google.golang.org/grpc"
)

//go:generate mockgen -source device_definitions_service.go -destination mocks/device_definitions_service_mock.go
type DeviceDefinitionsAPIService interface {
	GetDeviceDefinition(ctx context.Context, id string) (*pb.GetDeviceDefinitionItemResponse, error)
}

func NewDeviceDefinitionsAPIService(ddConn *grpc.ClientConn) DeviceDefinitionsAPIService {
	return &deviceDefinitionsAPIService{deviceDefinitionsConn: ddConn}
}

type deviceDefinitionsAPIService struct {
	deviceDefinitionsConn *grpc.ClientConn
}

func (dda *deviceDefinitionsAPIService) GetDeviceDefinition(ctx context.Context, id string) (*pb.GetDeviceDefinitionItemResponse, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("device definition id was empty - invalid")
	}
	definitionsClient := pb.NewDeviceDefinitionServiceClient(dda.deviceDefinitionsConn)

	def, err := definitionsClient.GetDeviceDefinitionByID(ctx, &pb.GetDeviceDefinitionRequest{
		Ids: []string{id},
	})
	if err != nil {
		return nil, err
	}

	return def.DeviceDefinitions[0], nil
}

func (dda *deviceDefinitionsAPIService) GetDeviceStyle(ctx context.Context, id string) (*pb.DeviceStyle, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("device style id was empty - invalid")
	}
	definitionsClient := pb.NewDeviceDefinitionServiceClient(dda.deviceDefinitionsConn)

	def, err := definitionsClient.GetDeviceStyleByID(ctx, &pb.GetDeviceStyleByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return def, nil
}
