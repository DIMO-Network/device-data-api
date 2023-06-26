package services

import (
	"context"
	"encoding/json"
	"fmt"

	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/volatiletech/sqlboiler/v4/boil"

	models "github.com/DIMO-Network/device-data-api/models"
	ddgrpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared/db"

	"github.com/rs/zerolog"
)

//go:generate mockgen -source vehicle_data_tracking_service.go -destination mocks/vehicle_data_tracking_service_mock.go
type VehicleDataTrackingService interface {
	GenerateVehicleDataTracking(ctx context.Context, userDeviceData models.UserDeviceDatum, deviceDefinition ddgrpc.GetDeviceDefinitionItemResponse, integration pb.UserDeviceIntegration) error
}

func NewVehicleDataTrackingService(db func() *db.ReaderWriter,
	log *zerolog.Logger) VehicleDataTrackingService {
	return &vehicleDataTrackingService{
		db:  db,
		log: log,
	}
}

type vehicleDataTrackingService struct {
	db          func() *db.ReaderWriter
	log         *zerolog.Logger
	memoryCache *gocache.Cache
}

func (v *vehicleDataTrackingService) GenerateVehicleDataTracking(ctx context.Context, userDeviceData models.UserDeviceDatum, deviceDefinition ddgrpc.GetDeviceDefinitionItemResponse, integration pb.UserDeviceIntegration) error {

	const CacheKey = "VehicleDataTrackingProperties"
	get, found := v.memoryCache.Get(CacheKey)

	eventAvailableProperties := make(map[string]string)
	if found {
		eventAvailableProperties = get.(map[string]string)
	} else {
		availableProperties, err := models.VehicleSignalsTrackingProperties().All(ctx, v.db().Reader)
		if err != nil {
			return err
		}
		for i := 0; i < len(availableProperties); i++ {
			eventAvailableProperties[availableProperties[i].Name] = availableProperties[i].ID
		}
		v.memoryCache.Set(CacheKey, eventAvailableProperties, 30*time.Minute)
	}

	var data map[string]interface{}
	err := json.Unmarshal(userDeviceData.Signals.JSON, &data)
	if err != nil {
		return nil
	}

	for key, value := range eventAvailableProperties {
		if _, ok := data[key]; ok {
			eventProperties := &models.VehicleSignalsTrackingEventsProperty{
				IntegrationID: integration.Id,
				DeviceMakeID:  deviceDefinition.Make.Id,
				PropertyID:    value,
				Year:          int(deviceDefinition.Type.Year),
				Model:         deviceDefinition.Type.Model,
				Count:         0,
			}

			if err := eventProperties.Upsert(ctx,
				v.db().Writer, true,
				[]string{models.VehicleSignalsTrackingEventsPropertyColumns.IntegrationID,
					models.VehicleSignalsTrackingEventsPropertyColumns.DeviceMakeID,
					models.VehicleSignalsTrackingEventsPropertyColumns.PropertyID},
				boil.Infer(), boil.Infer()); err != nil {
				return fmt.Errorf("error upserting VehicleDataTrackingEventsProperty: %w", err)
			}

			eventTracking, err := models.FindVehicleSignalsTrackingEventsMissingProperty(ctx, v.db().Writer, integration.Id, deviceDefinition.Make.Id, value)
			if err == nil {
				eventTracking.Count++
				_, err = eventTracking.Update(ctx, v.db().Writer, boil.Infer())
				if err != nil {
					v.log.Fatal()
				}
			}

		} else {
			eventMissingProperties := models.VehicleSignalsTrackingEventsMissingProperty{
				IntegrationID: integration.Id,
				DeviceMakeID:  deviceDefinition.Make.Id,
				PropertyID:    value,
				Year:          int(deviceDefinition.Type.Year),
				Model:         deviceDefinition.Type.Model,
				Count:         0,
			}

			if err := eventMissingProperties.Upsert(ctx,
				v.db().Writer, true,
				[]string{models.VehicleSignalsTrackingEventsMissingPropertyColumns.IntegrationID,
					models.VehicleSignalsTrackingEventsMissingPropertyColumns.DeviceMakeID,
					models.VehicleSignalsTrackingEventsMissingPropertyColumns.PropertyID},
				boil.Infer(), boil.Infer()); err != nil {
				return fmt.Errorf("error upserting eventMissingProperties: %w", err)
			}

			eventTracking, err := models.FindVehicleSignalsTrackingEventsMissingProperty(ctx, v.db().Writer, integration.Id, deviceDefinition.Make.Id, value)
			if err == nil {
				eventTracking.Count++
				_, err = eventTracking.Update(ctx, v.db().Writer, boil.Infer())
				if err != nil {
					v.log.Fatal()
				}
			}
		}
	}

	return nil
}
