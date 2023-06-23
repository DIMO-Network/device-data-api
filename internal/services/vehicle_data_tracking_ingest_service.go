package services

import (
	"context"
	"fmt"
	"time"

	"github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/shared/db"
	"github.com/lovoo/goka"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

type VehicleDataTrackingIngestService struct {
	db           func() *db.ReaderWriter
	log          *zerolog.Logger
	eventService EventService
	deviceDefSvc DeviceDefinitionsAPIService
	integrations []*grpc.Integration
	autoPiSvc    AutoPiAPIService
	deviceSvc    DeviceAPIService
	memoryCache  *gocache.Cache
}

func NewVehicleDataTrackingIngestService(db func() *db.ReaderWriter, log *zerolog.Logger, eventService EventService, ddSvc DeviceDefinitionsAPIService, autoPiSvc AutoPiAPIService, deviceSvc DeviceAPIService) *VehicleDataTrackingIngestService {
	// Cache the list of integrations.
	integrations, err := ddSvc.GetIntegrations(context.Background())
	if err != nil {
		log.Fatal().Err(err).Str("func", "NewDeviceStatusIngestService").Msg("Couldn't retrieve global integration list.")
	}
	c := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids

	return &VehicleDataTrackingIngestService{
		db:           db,
		log:          log,
		deviceDefSvc: ddSvc,
		eventService: eventService,
		integrations: integrations,
		autoPiSvc:    autoPiSvc,
		deviceSvc:    deviceSvc,
		memoryCache:  c,
	}
}

// ProcessDeviceStatusMessages works on channel stream of messages from watermill kafka consumer
func (i *VehicleDataTrackingIngestService) ProcessDeviceStatusMessages(ctx goka.Context, msg interface{}) {
	if err := i.processMessage(ctx, msg.(*DeviceStatusEvent)); err != nil {
		i.log.Err(err).Msg("Error processing device status message.")
	}
}

func (i *VehicleDataTrackingIngestService) processMessage(ctx goka.Context, event *DeviceStatusEvent) error {
	if event.Type != deviceStatusEventType {
		return fmt.Errorf("received vehicle status event with unexpected type %s", event.Type)
	}

	return nil
}
