package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"google.golang.org/grpc"
)

//go:generate mockgen -source device_definitions_service.go -destination mocks/device_definitions_service_mock.go
type DeviceDefinitionsAPIService interface {
	GetDeviceDefinitionByID(ctx context.Context, id string) (*pb.GetDeviceDefinitionItemResponse, error)
	GetIntegrations(ctx context.Context) ([]*pb.Integration, error)
	GetDeviceDefinitionsByIDs(ctx context.Context, ids []string) ([]*pb.GetDeviceDefinitionItemResponse, error)
}

func NewDeviceDefinitionsAPIService(ddConn *grpc.ClientConn) DeviceDefinitionsAPIService {
	return &deviceDefinitionsAPIService{deviceDefinitionsConn: ddConn}
}

type deviceDefinitionsAPIService struct {
	deviceDefinitionsConn *grpc.ClientConn
}

// GetDeviceDefinitionsByIDs calls device definitions api via GRPC to get the definition. idea for testing: http://www.inanzzz.com/index.php/post/w9qr/unit-testing-golang-grpc-client-and-server-application-with-bufconn-package
// if not found or other error from server, the error contains the grpc status code that can be interpreted for different conditions. example in api.GrpcErrorToFiber
func (dda *deviceDefinitionsAPIService) GetDeviceDefinitionsByIDs(ctx context.Context, ids []string) ([]*pb.GetDeviceDefinitionItemResponse, error) {

	if len(ids) == 0 {
		return nil, errors.New("Device Definition Ids is required")
	}

	definitionsClient := pb.NewDeviceDefinitionServiceClient(dda.deviceDefinitionsConn)

	definitions, err2 := definitionsClient.GetDeviceDefinitionByID(ctx, &pb.GetDeviceDefinitionRequest{
		Ids: ids,
	})

	if err2 != nil {
		return nil, err2
	}

	return definitions.GetDeviceDefinitions(), nil
}

// GetIntegrations calls device definitions integrations api via GRPC to get the definition. idea for testing: http://www.inanzzz.com/index.php/post/w9qr/unit-testing-golang-grpc-client-and-server-application-with-bufconn-package
func (dda *deviceDefinitionsAPIService) GetIntegrations(ctx context.Context) ([]*pb.Integration, error) {
	definitionsClient := pb.NewDeviceDefinitionServiceClient(dda.deviceDefinitionsConn)

	definitions, err := definitionsClient.GetIntegrations(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return definitions.GetIntegrations(), nil
}

func (dda *deviceDefinitionsAPIService) GetDeviceDefinitionByID(ctx context.Context, id string) (*pb.GetDeviceDefinitionItemResponse, error) {
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

func GetActualDeviceDefinitionMetadataValues(dd *pb.GetDeviceDefinitionItemResponse, deviceStyleID *string) *DeviceDefinitionRange {

	var fuelTankCapGal, mpg, mpgHwy float64 = 0, 0, 0

	var metadata []*pb.DeviceTypeAttribute

	if deviceStyleID != nil {
		for _, style := range dd.DeviceStyles {
			if style.Id == *deviceStyleID {
				metadata = style.DeviceAttributes
				break
			}
		}
	}

	if len(metadata) == 0 && dd != nil && dd.DeviceAttributes != nil {
		metadata = dd.DeviceAttributes
	}

	for _, attr := range metadata {
		switch DeviceAttributeType(attr.Name) {
		case FuelTankCapacityGal:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				fuelTankCapGal = v
			}
		case MPG:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				mpg = v
			}
		case MpgHighway:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				mpgHwy = v
			}
		}
	}

	return &DeviceDefinitionRange{
		FuelTankCapGal: fuelTankCapGal,
		Mpg:            mpg,
		MpgHwy:         mpgHwy,
	}
}

type DeviceAttributeType string

const (
	MPG                 DeviceAttributeType = "mpg"
	FuelTankCapacityGal DeviceAttributeType = "fuel_tank_capacity_gal"
	MpgHighway          DeviceAttributeType = "mpg_highway"
)

type DeviceDefinitionRange struct {
	FuelTankCapGal float64 `json:"fuel_tank_capacity_gal"`
	Mpg            float64 `json:"mpg"`
	MpgHwy         float64 `json:"mpg_highway"`
}
