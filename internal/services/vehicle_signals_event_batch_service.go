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
	GenerateVehicleDataTracking(ctx context.Context) error
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

func (v *vehicleSignalsEventBatchService) GenerateVehicleDataTracking(ctx context.Context) error {

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

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	fromTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	toTime := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 0, yesterday.Location())

	deviceDataEvents, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GTE(fromTime),
		models.UserDeviceDatumWhere.UpdatedAt.LTE(toTime),
	).All(ctx, v.db().Reader)
	if err != nil {
		return err
	}

	dateKey := fromTime.Format("20060102")

	for _, item := range deviceDataEvents {

		device := &pb.UserDevice{}
		get, found := v.memoryCache.Get(item.UserDeviceID + "_" + item.IntegrationID.String)
		if found {
			device = get.(*pb.UserDevice)
		} else {
			// todo problem is here - does not get back the device.DeviceDefinitionId
			// problem is with v.deviceSvc not being instantiated correctly
			device, err := v.deviceSvc.GetUserDevice(ctx, item.UserDeviceID)
			if err != nil {
				v.log.Err(err).Msgf("failed to find device %s", item.UserDeviceID)
				continue
			}
			// Validate integration Id
			v.memoryCache.Set(item.UserDeviceID+"_"+item.IntegrationID.String, device, 30*time.Minute)
		}

		v.log.Info().Msgf("DeviceID %s, UserDeviceID %s, DeviceDefinitionId %s, FromCache %v", device.Id, item.UserDeviceID, device.DeviceDefinitionId, found)

		deviceDefinition, err := v.deviceDefSvc.GetDeviceDefinitionByID(ctx, device.DeviceDefinitionId)
		if err != nil {
			v.log.Err(err).Msgf("deviceDefSvc error getting definition id: %s", device.DeviceDefinitionId)
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

		for key, value := range eventAvailableProperties {

			integrationID := item.IntegrationID.String
			deviceMakeID := deviceDefinition.Make.Id
			model := deviceDefinition.Type.Model
			year := int(deviceDefinition.Type.Year)

			event, err := models.ReportVehicleSignalsEvents(
				models.ReportVehicleSignalsEventsPropertyWhere.DateID.EQ(dateKey),
				models.ReportVehicleSignalsEventsPropertyWhere.IntegrationID.EQ(integrationID),
				models.ReportVehicleSignalsEventsPropertyWhere.DeviceMakeID.EQ(deviceMakeID),
				models.ReportVehicleSignalsEventsPropertyWhere.PropertyID.EQ(key),
				models.ReportVehicleSignalsEventsPropertyWhere.Model.EQ(model),
				models.ReportVehicleSignalsEventsPropertyWhere.Year.EQ(year),
			).One(ctx, v.db().Reader)

			if err != nil {
				if err != sql.ErrNoRows {
					v.log.Err(err).Msg("failed to find report vehicle signals")
					continue
				}
			}

			if event == nil {
				event = &models.ReportVehicleSignalsEvent{
					DateID:             dateKey,
					IntegrationID:      item.IntegrationID.String,
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
				models.ReportVehicleSignalsEventColumns.DateID,
				models.ReportVehicleSignalsEventColumns.IntegrationID,
				models.ReportVehicleSignalsEventColumns.DeviceMakeID,
				models.ReportVehicleSignalsEventColumns.PropertyID,
				models.ReportVehicleSignalsEventColumns.Model,
				models.ReportVehicleSignalsEventColumns.Year,
			}

			if err := event.Upsert(ctx, v.db().Writer, true, reportVehicleSignalsEventPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
				v.log.Err(err).Msgf("error inserting report event : %s %s %s %s", event.DateID, event.IntegrationID, event.DeviceMakeID, event.PropertyID)
			}

			if _, ok := data[key]; ok {
				eventProperties, err := models.ReportVehicleSignalsEventsProperties(
					models.ReportVehicleSignalsEventsPropertyWhere.DateID.EQ(dateKey),
					models.ReportVehicleSignalsEventsPropertyWhere.IntegrationID.EQ(integrationID),
					models.ReportVehicleSignalsEventsPropertyWhere.DeviceMakeID.EQ(deviceMakeID),
					models.ReportVehicleSignalsEventsPropertyWhere.PropertyID.EQ(key),
					models.ReportVehicleSignalsEventsPropertyWhere.Model.EQ(model),
					models.ReportVehicleSignalsEventsPropertyWhere.Year.EQ(year),
				).One(ctx, v.db().Reader)

				if err != nil {
					if err != sql.ErrNoRows {
						v.log.Err(err).Msg("failed to find report vehicle signals")
						continue
					}
				}

				if eventProperties == nil {
					eventProperties = &models.ReportVehicleSignalsEventsProperty{
						DateID:             dateKey,
						IntegrationID:      item.IntegrationID.String,
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
					models.ReportVehicleSignalsEventsPropertyColumns.DateID,
					models.ReportVehicleSignalsEventsPropertyColumns.IntegrationID,
					models.ReportVehicleSignalsEventsPropertyColumns.DeviceMakeID,
					models.ReportVehicleSignalsEventsPropertyColumns.PropertyID,
					models.ReportVehicleSignalsEventsPropertyColumns.Model,
					models.ReportVehicleSignalsEventsPropertyColumns.Year,
				}

				if err := eventProperties.Upsert(ctx, v.db().Writer, true, reportVehicleSignalsPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
					v.log.Err(err).Msgf("error inserting report properties : %s %s %s %s", eventProperties.DateID, eventProperties.IntegrationID, eventProperties.DeviceMakeID, eventProperties.PropertyID)
				}
			}
		}

	}

	return nil
}
