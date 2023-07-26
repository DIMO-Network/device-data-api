package services

import (
	"context"
	"time"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"

	models "github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared/db"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

//go:generate mockgen -source vehicle_data_tracking_service.go -destination mocks/vehicle_data_tracking_service_mock.go
type VehicleSignalsEventBatchService interface {
	GenerateVehicleDataTracking(ctx context.Context, dateKey string, fromTime time.Time) error
}

func NewVehicleSignalsEventBatchService(db func() *db.ReaderWriter,
	log *zerolog.Logger, deviceDefSvc DeviceDefinitionsAPIService, deviceSvc DeviceAPIService,
	vehicleSignalsEventPropertyService VehicleSignalsEventPropertyService, vehicleSignalsEventDeviceUserService VehicleSignalsEventDeviceUserService) VehicleSignalsEventBatchService {
	cache := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids
	return &vehicleSignalsEventBatchService{
		db:                                   db,
		log:                                  log,
		deviceDefSvc:                         deviceDefSvc,
		deviceSvc:                            deviceSvc,
		memoryCache:                          cache,
		vehicleSignalsEventPropertyService:   vehicleSignalsEventPropertyService,
		vehicleSignalsEventDeviceUserService: vehicleSignalsEventDeviceUserService,
	}
}

type vehicleSignalsEventBatchService struct {
	db                                   func() *db.ReaderWriter
	log                                  *zerolog.Logger
	memoryCache                          *gocache.Cache
	deviceDefSvc                         DeviceDefinitionsAPIService
	deviceSvc                            DeviceAPIService
	vehicleSignalsEventPropertyService   VehicleSignalsEventPropertyService
	vehicleSignalsEventDeviceUserService VehicleSignalsEventDeviceUserService
}

func (v *vehicleSignalsEventBatchService) GenerateVehicleDataTracking(ctx context.Context, dateKey string, fromTime time.Time) error {

	const CacheKey = "VehicleDataTrackingProperties"
	get, found := v.memoryCache.Get(CacheKey)

	eventAvailableProperties := make(map[string]string)
	if found {
		eventAvailableProperties = get.(map[string]string)
	} else {
		availableProperties, err := models.VehicleSignalsAvailableProperties().All(ctx, v.db().Reader)
		if err != nil {
			return err
		}
		for i := 0; i < len(availableProperties); i++ {
			eventAvailableProperties[availableProperties[i].Name] = availableProperties[i].ID
		}
		v.memoryCache.Set(CacheKey, eventAvailableProperties, 30*time.Minute)
	}

	v.log.Info().Msgf("Available properties: %d", len(eventAvailableProperties))

	deviceDataEvents, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GTE(fromTime),
		//models.UserDeviceDatumWhere.UpdatedAt.LTE(toTime),
	).All(ctx, v.db().Reader)
	if err != nil {
		return err
	}

	v.log.Info().Msgf("snapshot based on userDeviceData records: %d, where updated_at > %s", len(deviceDataEvents), fromTime.Format(time.RFC822))

	for _, item := range deviceDataEvents {

		device := &pb.UserDevice{}
		cachedUD, foundCached := v.memoryCache.Get(item.UserDeviceID + "_" + item.IntegrationID)
		if foundCached {
			device = cachedUD.(*pb.UserDevice)
		} else {
			device, err = v.deviceSvc.GetUserDevice(ctx, item.UserDeviceID)
			if err != nil {
				v.log.Err(err).Msgf("failed to find device %s", item.UserDeviceID)
				continue
			}
			//v.log.Info().Msgf("found UD: \n %+v", device)
			// Validate integration Id
			v.memoryCache.Set(item.UserDeviceID+"_"+item.IntegrationID, device, 30*time.Minute)
		}

		deviceDefinition, err := v.deviceDefSvc.GetDeviceDefinitionByID(ctx, device.DeviceDefinitionId)
		if err != nil {
			v.log.Err(err).Msgf("(%s) deviceDefSvc error getting definition id: %s", device.Id, device.DeviceDefinitionId)
			continue
		}

		if deviceDefinition == nil {
			v.log.Err(err).Msgf("device definition with id %s not found", device.DeviceDefinitionId)
			continue
		}

		err = v.vehicleSignalsEventPropertyService.GenerateData(ctx, dateKey, item.IntegrationID, item, deviceDefinition, eventAvailableProperties)
		if err != nil {
			v.log.Err(err).Msgf("(%s) generate event property error: %s", device.Id, device.DeviceDefinitionId)
			continue
		}

		err = v.vehicleSignalsEventDeviceUserService.GenerateData(ctx, dateKey, item.IntegrationID, device.PowerTrainType)
		if err != nil {
			v.log.Err(err).Msgf("(%s) generate user device error: %s", device.Id, device.DeviceDefinitionId)
			continue
		}
	}

	return nil
}
