package services

import (
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"

	"github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared/db"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

//go:generate mockgen -source vehicle_signals_event_user_device_service.go -destination mocks/vehicle_signals_event_user_device_service_mock.go
type vehicleSignalsEventDeviceUserService struct {
	db          func() *db.ReaderWriter
	log         *zerolog.Logger
	memoryCache *gocache.Cache
}

type VehicleSignalsEventDeviceUserService interface {
	GenerateData(ctx context.Context, dateKey string, integrationID string, powerTrainType string) error
}

func NewVehicleSignalsEventDeviceUserService(db func() *db.ReaderWriter, log *zerolog.Logger) VehicleSignalsEventDeviceUserService {
	cache := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids
	return &vehicleSignalsEventDeviceUserService{
		db:          db,
		log:         log,
		memoryCache: cache,
	}
}

func (v *vehicleSignalsEventDeviceUserService) GenerateData(ctx context.Context, dateKey string, integrationID string, powerTrainType string) error {

	userDeviceEvent, err := models.ReportVehicleSignalsEventsUserDevices(
		models.ReportVehicleSignalsEventsUserDeviceWhere.DateID.EQ(dateKey),
		models.ReportVehicleSignalsEventsUserDeviceWhere.IntegrationID.EQ(integrationID),
		models.ReportVehicleSignalsEventsUserDeviceWhere.PowerTrainType.EQ(powerTrainType),
	).One(ctx, v.db().Reader)

	if err != nil {
		if err != sql.ErrNoRows {
			v.log.Err(err).Msg("failed to find report vehicle signals")
			return err
		}
	}

	if userDeviceEvent == nil {
		userDeviceEvent = &models.ReportVehicleSignalsEventsUserDevice{
			DateID:         dateKey,
			IntegrationID:  integrationID,
			PowerTrainType: powerTrainType,
			Count:          1,
		}
	} else {
		userDeviceEvent.Count++
	}

	var reportVehicleSignalsPrimaryKeyColumns = []string{
		models.ReportVehicleSignalsEventsUserDeviceColumns.DateID,
		models.ReportVehicleSignalsEventsUserDeviceColumns.IntegrationID,
		models.ReportVehicleSignalsEventsUserDeviceColumns.PowerTrainType,
	}

	if err := userDeviceEvent.Upsert(ctx, v.db().Writer, true, reportVehicleSignalsPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
		v.log.Err(err).Msgf("error inserting report user device : %s %s", userDeviceEvent.DateID, userDeviceEvent.IntegrationID)
	}

	return nil
}
