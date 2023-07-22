package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"

	models "github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared/db"
	gocache "github.com/patrickmn/go-cache"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/rs/zerolog"
)

//go:generate mockgen -source vehicle_data_tracking_service.go -destination mocks/vehicle_data_tracking_service_mock.go
type VehicleSignalsEventBatchService interface {
	GenerateVehicleDataTracking(ctx context.Context, dateKey string, fromTime time.Time) error
}

func NewVehicleSignalsEventBatchService(db func() *db.ReaderWriter,
	log *zerolog.Logger, deviceDefSvc DeviceDefinitionsAPIService, deviceSvc DeviceAPIService) VehicleSignalsEventBatchService {
	cache := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids
	return &vehicleSignalsEventBatchService{
		db:           db,
		log:          log,
		deviceDefSvc: deviceDefSvc,
		deviceSvc:    deviceSvc,
		memoryCache:  cache,
	}
}

type vehicleSignalsEventBatchService struct {
	db           func() *db.ReaderWriter
	log          *zerolog.Logger
	memoryCache  *gocache.Cache
	deviceDefSvc DeviceDefinitionsAPIService
	deviceSvc    DeviceAPIService
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

		var data map[string]interface{}
		err = json.Unmarshal(item.Signals.JSON, &data)
		if err != nil {
			continue
		}

		v.log.Info().Msgf("Available properties: %d", len(eventAvailableProperties))

		for key, value := range eventAvailableProperties {

			integrationID := item.IntegrationID
			deviceMakeID := deviceDefinition.Make.Id
			model := deviceDefinition.Type.Model
			year := int(deviceDefinition.Type.Year)

			event, err := models.ReportVehicleSignalsEventsAlls(
				models.ReportVehicleSignalsEventsAllWhere.DateID.EQ(dateKey),
				models.ReportVehicleSignalsEventsAllWhere.IntegrationID.EQ(integrationID),
				models.ReportVehicleSignalsEventsAllWhere.DeviceMakeID.EQ(deviceMakeID),
				models.ReportVehicleSignalsEventsAllWhere.PropertyID.EQ(key),
				models.ReportVehicleSignalsEventsAllWhere.Model.EQ(model),
				models.ReportVehicleSignalsEventsAllWhere.Year.EQ(year),
			).One(ctx, v.db().Reader)

			if err != nil {
				if err != sql.ErrNoRows {
					v.log.Err(err).Msg("failed to find report vehicle signals")
					continue
				}
			}

			if event == nil {
				event = &models.ReportVehicleSignalsEventsAll{
					DateID:             dateKey,
					IntegrationID:      item.IntegrationID,
					DeviceMakeID:       deviceDefinition.Make.Id,
					PropertyID:         value,
					Year:               int(deviceDefinition.Type.Year),
					Model:              deviceDefinition.Type.Model,
					DeviceDefinitionID: deviceDefinition.DeviceDefinitionId,
					DeviceMake:         deviceDefinition.Make.Name,
					Count:              0,
				}
			} else {
				event.Count++
			}

			var reportVehicleSignalsEventPrimaryKeyColumns = []string{
				models.ReportVehicleSignalsEventsAllColumns.DateID,
				models.ReportVehicleSignalsEventsAllColumns.IntegrationID,
				models.ReportVehicleSignalsEventsAllColumns.DeviceMakeID,
				models.ReportVehicleSignalsEventsAllColumns.PropertyID,
				models.ReportVehicleSignalsEventsAllColumns.Model,
				models.ReportVehicleSignalsEventsAllColumns.Year,
			}

			if err := event.Upsert(ctx, v.db().Writer, true, reportVehicleSignalsEventPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
				v.log.Err(err).Msgf("error inserting report event : %s %s %s %s", event.DateID, event.IntegrationID, event.DeviceMakeID, event.PropertyID)
			}

			if _, ok := data[key]; ok {
				eventProperties, err := models.ReportVehicleSignalsEventsTrackings(
					models.ReportVehicleSignalsEventsTrackingWhere.DateID.EQ(dateKey),
					models.ReportVehicleSignalsEventsTrackingWhere.IntegrationID.EQ(integrationID),
					models.ReportVehicleSignalsEventsTrackingWhere.DeviceMakeID.EQ(deviceMakeID),
					models.ReportVehicleSignalsEventsTrackingWhere.PropertyID.EQ(key),
					models.ReportVehicleSignalsEventsTrackingWhere.Model.EQ(model),
					models.ReportVehicleSignalsEventsTrackingWhere.Year.EQ(year),
				).One(ctx, v.db().Reader)

				if err != nil {
					if err != sql.ErrNoRows {
						v.log.Err(err).Msg("failed to find report vehicle signals")
						continue
					}
				}

				if eventProperties == nil {
					eventProperties = &models.ReportVehicleSignalsEventsTracking{
						DateID:             dateKey,
						IntegrationID:      item.IntegrationID,
						DeviceMakeID:       deviceDefinition.Make.Id,
						PropertyID:         value,
						Year:               int(deviceDefinition.Type.Year),
						Model:              deviceDefinition.Type.Model,
						DeviceDefinitionID: deviceDefinition.DeviceDefinitionId,
						DeviceMake:         deviceDefinition.Make.Name,
						Count:              0,
					}
				} else {
					eventProperties.Count++
				}

				var reportVehicleSignalsPrimaryKeyColumns = []string{
					models.ReportVehicleSignalsEventsTrackingColumns.DateID,
					models.ReportVehicleSignalsEventsTrackingColumns.IntegrationID,
					models.ReportVehicleSignalsEventsTrackingColumns.DeviceMakeID,
					models.ReportVehicleSignalsEventsTrackingColumns.PropertyID,
					models.ReportVehicleSignalsEventsTrackingColumns.Model,
					models.ReportVehicleSignalsEventsTrackingColumns.Year,
				}

				if err := eventProperties.Upsert(ctx, v.db().Writer, true, reportVehicleSignalsPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
					v.log.Err(err).Msgf("error inserting report properties : %s %s %s %s", eventProperties.DateID, eventProperties.IntegrationID, eventProperties.DeviceMakeID, eventProperties.PropertyID)
				}
			}
		}

	}

	return nil
}
