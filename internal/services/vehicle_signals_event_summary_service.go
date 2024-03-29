package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared/db"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

type vehicleSignalsEventSummaryService struct {
	db          func() *db.ReaderWriter
	log         *zerolog.Logger
	memoryCache *gocache.Cache
}

type VehicleSignalsEventSummaryService interface {
	GenerateData(ctx context.Context, dateKey string, integrationID string, powerTrainType string, deviceDefinitionCount int) error
}

func NewVehicleSignalsEventSummaryService(db func() *db.ReaderWriter, log *zerolog.Logger) VehicleSignalsEventSummaryService {
	cache := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids
	return &vehicleSignalsEventSummaryService{
		db:          db,
		log:         log,
		memoryCache: cache,
	}
}

func (v *vehicleSignalsEventSummaryService) GenerateData(ctx context.Context, dateKey string, integrationID string, powerTrainType string, deviceDefinitionCount int) error {

	userDeviceEvent, err := models.ReportVehicleSignalsEventsSummaries(
		models.ReportVehicleSignalsEventsSummaryWhere.DateID.EQ(dateKey),
		models.ReportVehicleSignalsEventsSummaryWhere.IntegrationID.EQ(integrationID),
		models.ReportVehicleSignalsEventsSummaryWhere.PowerTrainType.EQ(powerTrainType),
	).One(ctx, v.db().Reader)

	if err != nil {
		if err != sql.ErrNoRows {
			v.log.Err(err).Msg("failed to find report vehicle signals")
			return err
		}
	}

	if userDeviceEvent == nil {
		userDeviceEvent = &models.ReportVehicleSignalsEventsSummary{
			DateID:                dateKey,
			IntegrationID:         integrationID,
			PowerTrainType:        powerTrainType,
			DeviceDefinitionCount: deviceDefinitionCount,
			Count:                 1,
		}
	} else {
		userDeviceEvent.Count++
		userDeviceEvent.DeviceDefinitionCount = deviceDefinitionCount
	}

	var reportVehicleSignalsPrimaryKeyColumns = []string{
		models.ReportVehicleSignalsEventsSummaryColumns.DateID,
		models.ReportVehicleSignalsEventsSummaryColumns.IntegrationID,
		models.ReportVehicleSignalsEventsSummaryColumns.PowerTrainType,
	}

	if err := userDeviceEvent.Upsert(ctx, v.db().Writer, true, reportVehicleSignalsPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
		v.log.Err(err).Msgf("error inserting report user device : %s %s", userDeviceEvent.DateID, userDeviceEvent.IntegrationID)
	}

	return nil
}
