package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"time"

	"github.com/DIMO-Network/device-data-api/models"
	pb "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	gocache "github.com/patrickmn/go-cache"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
)

type vehicleSignalsEventPropertyService struct {
	db          func() *db.ReaderWriter
	log         *zerolog.Logger
	memoryCache *gocache.Cache
}

type VehicleSignalsEventPropertyService interface {
	GenerateData(ctx context.Context, dateKey string, integrationID string, ud *models.UserDeviceDatum, deviceDefinition *pb.GetDeviceDefinitionItemResponse, eventAvailableProperties map[string]models.VehicleSignalsAvailableProperty) error
}

func NewVehicleSignalsEventPropertyService(db func() *db.ReaderWriter, log *zerolog.Logger) VehicleSignalsEventPropertyService {
	cache := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids
	return &vehicleSignalsEventPropertyService{
		db:          db,
		log:         log,
		memoryCache: cache,
	}
}

func (v *vehicleSignalsEventPropertyService) GenerateData(ctx context.Context, dateKey string, integrationID string, ud *models.UserDeviceDatum, deviceDefinition *pb.GetDeviceDefinitionItemResponse, eventAvailableProperties map[string]models.VehicleSignalsAvailableProperty) error {

	var data map[string]interface{}
	err := json.Unmarshal(ud.Signals.JSON, &data)
	if err != nil {
		return err
	}

	for key, value := range eventAvailableProperties {

		// Validate if the property should be considered
		if !v.shouldValueBeConsidered(key, value, data) {
			continue
		}

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
				IntegrationID:      integrationID,
				DeviceMakeID:       deviceDefinition.Make.Id,
				PropertyID:         value.ID,
				Year:               int(deviceDefinition.Type.Year),
				Model:              deviceDefinition.Type.Model,
				DeviceDefinitionID: deviceDefinition.DeviceDefinitionId,
				DeviceMake:         deviceDefinition.Make.Name,
				Count:              1,
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
					IntegrationID:      integrationID,
					DeviceMakeID:       deviceDefinition.Make.Id,
					PropertyID:         value.ID,
					Year:               int(deviceDefinition.Type.Year),
					Model:              deviceDefinition.Type.Model,
					DeviceDefinitionID: deviceDefinition.DeviceDefinitionId,
					DeviceMake:         deviceDefinition.Make.Name,
					Count:              1,
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

	return nil
}

func (v *vehicleSignalsEventPropertyService) shouldValueBeConsidered(key string, value models.VehicleSignalsAvailableProperty, data map[string]interface{}) bool {
	if _, ok := data[key]; ok {
		if value.ValidMinLength.Valid && value.ValidMinLength.Int > 0 {
			object, ok := data[key].(map[string]interface{})
			if !ok {
				return false
			}

			objectValue, ok := object["value"].(string)
			if !ok {
				return false
			}

			if len(objectValue) < value.ValidMinLength.Int {
				return false
			}
		}

		return true
	}

	return false
}
