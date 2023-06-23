package services

import (
	"context"

	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
)

//go:generate mockgen -source vehicle_data_tracking_service.go -destination mocks/vehicle_data_tracking_service_mock.go
type VehicleDataTrackingService interface {
	GenerateVehicleDataTracking(ctx context.Context) error
}

func NewVehicleDataTrackingService(db func() *db.ReaderWriter,
	log *zerolog.Logger,
	ddSvc DeviceDefinitionsAPIService,
	deviceSvc DeviceAPIService) VehicleDataTrackingService {
	return &vehicleDataTrackingService{
		db:           db,
		log:          log,
		deviceDefSvc: ddSvc,
		deviceSvc:    deviceSvc,
	}
}

type vehicleDataTrackingService struct {
	db           func() *db.ReaderWriter
	log          *zerolog.Logger
	deviceDefSvc DeviceDefinitionsAPIService
	deviceSvc    DeviceAPIService
}

func (v *vehicleDataTrackingService) GenerateVehicleDataTracking(ctx context.Context) error {
	return nil
}
