package main

import (
	"context"
	"flag"

	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/shared/db"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"
)

type vehicleSignalsBatchCmd struct {
	db           func() *db.ReaderWriter
	logger       zerolog.Logger
	deviceDefSvc services.DeviceDefinitionsAPIService
	deviceSvc    services.DeviceAPIService
}

func (*vehicleSignalsBatchCmd) Name() string {
	return "generate-report-vehicle-signals"
}
func (*vehicleSignalsBatchCmd) Synopsis() string {
	return "generate vehicle signals events report by date for Data Dashboard"
}
func (*vehicleSignalsBatchCmd) Usage() string {
	return `generate-report-vehicle-signals`
}

// nolint
func (p *vehicleSignalsBatchCmd) SetFlags(f *flag.FlagSet) {

}

func (p *vehicleSignalsBatchCmd) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	p.logger.Info().Msg("running batch report for vehicle signals")

	vehicleSignalsEventPropertySrv := services.NewVehicleSignalsEventPropertyService(p.db, &p.logger)
	vehicleSignalsEventSummarySrv := services.NewVehicleSignalsEventSummaryService(p.db, &p.logger)

	batchSrv := services.NewVehicleSignalsEventBatchService(p.db, &p.logger, p.deviceDefSvc, p.deviceSvc, vehicleSignalsEventPropertySrv, vehicleSignalsEventSummarySrv)
	vehicleSignalJobSrv := services.NewVehicleSignalJobService(p.db, &p.logger)
	jobContext, err := vehicleSignalJobSrv.GetJobContext(ctx)

	if err != nil {
		p.logger.Error().Err(err).Msg("Error job context")
	}

	p.logger.Log().Msgf("Execute : %v", jobContext.Execute)

	if jobContext.Execute {
		err := batchSrv.GenerateVehicleDataTracking(ctx, jobContext.DateKey, jobContext.FromTime)
		if err != nil {
			p.logger.Fatal().Err(err).Msg("Error running vehicle signals event batch service")
		}
	}

	return subcommands.ExitSuccess
}
