package main

import (
	"context"
	"flag"

	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/shared/db"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"
)

type vehicleSignalsEventBatchServiceCmd struct {
	db           func() *db.ReaderWriter
	logger       zerolog.Logger
	deviceDefSvc services.DeviceDefinitionsAPIService
	deviceSvc    services.DeviceAPIService
}

func (*vehicleSignalsEventBatchServiceCmd) Name() string {
	return "generate-report-vehicle-signals-event"
}
func (*vehicleSignalsEventBatchServiceCmd) Synopsis() string {
	return "generate vehicle signals events report by date"
}
func (*vehicleSignalsEventBatchServiceCmd) Usage() string {
	return `generate-report-vehicle-signals-event`
}

// nolint
func (p *vehicleSignalsEventBatchServiceCmd) SetFlags(f *flag.FlagSet) {

}

func (p *vehicleSignalsEventBatchServiceCmd) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	batchSrv := services.NewVehicleSignalsEventBatchService(p.db, &p.logger, p.deviceDefSvc, p.deviceSvc)
	err := batchSrv.GenerateVehicleDataTracking(ctx)
	if err != nil {
		p.logger.Fatal().Err(err).Msg("Error running vehicle signals event batch service")
	}
	return subcommands.ExitSuccess
}