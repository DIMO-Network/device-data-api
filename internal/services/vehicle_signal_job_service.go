package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/DIMO-Network/device-data-api/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	internalmodel "github.com/DIMO-Network/device-data-api/internal/models"
	"github.com/DIMO-Network/shared/db"
	"github.com/rs/zerolog"
)

//go:generate mockgen -source vehicle_signal_job_service.go -destination mocks/vehicle_signal_job_service_mock.go
type vehicleSignalJobService struct {
	db  func() *db.ReaderWriter
	log *zerolog.Logger
}

type VehicleSignalJobService interface {
	GetJobContext(ctx context.Context) (*internalmodel.SignalJobContext, error)
}

func NewVehicleSignalJobService(db func() *db.ReaderWriter, log *zerolog.Logger) VehicleSignalJobService {
	return &vehicleSignalJobService{
		db:  db,
		log: log,
	}
}

func (v vehicleSignalJobService) GetJobContext(ctx context.Context) (*internalmodel.SignalJobContext, error) {
	queryMods := []qm.QueryMod{
		qm.OrderBy("created_at DESC"),
		qm.Limit(1),
	}

	vehicleSignalsJob, err := models.VehicleSignalsJobs(queryMods...).One(ctx, v.db().Reader)
	if err != nil {
		if err != sql.ErrNoRows {
			v.log.Err(err).Msg("failed to find signal job")
			return nil, err
		}
	}

	now := time.Now()
	startDate := now.AddDate(0, 0, -7)
	endDate := now
	fromTime := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())
	dateKey := fromTime.Format("20060102")

	v.log.Info().Msgf("Current Date:  %s", now)
	v.log.Info().Msgf("From Date:  %s", startDate)
	v.log.Info().Msgf("End Date:  %s", endDate)
	v.log.Info().Msgf("Date format:  %s", dateKey)

	if vehicleSignalsJob != nil {
		daysDifference := vehicleSignalsJob.EndDate.Sub(now).Hours() / 24
		v.log.Err(err).Msgf("Day Difference %f", daysDifference)
		if daysDifference > 7 {
			vehicleSignalsJob = &models.VehicleSignalsJob{
				VehicleSignalsJobID: dateKey,
				StartDate:           startDate,
				EndDate:             endDate,
			}
		} else {

			return &internalmodel.SignalJobContext{Execute: false}, nil
		}
	}

	if vehicleSignalsJob == nil {
		vehicleSignalsJob = &models.VehicleSignalsJob{
			VehicleSignalsJobID: dateKey,
			StartDate:           startDate,
			EndDate:             endDate,
		}
	}

	err = vehicleSignalsJob.Insert(ctx, v.db().Writer, boil.Infer())
	if err != nil {
		return nil, err
	}

	return &internalmodel.SignalJobContext{Execute: true, DateKey: dateKey, FromTime: fromTime}, nil
}
