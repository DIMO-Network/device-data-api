package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"

	"github.com/tidwall/gjson"

	"github.com/DIMO-Network/device-data-api/internal/appmetrics"
	"github.com/DIMO-Network/device-data-api/internal/constants"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/shared/db"
	"github.com/lovoo/goka"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const (
	deviceStatusEventType = "zone.dimo.device.status.update"
	odometerCooldown      = time.Hour
)

type DeviceStatusIngestService struct {
	db           func() *db.ReaderWriter
	log          *zerolog.Logger
	eventService EventService
	deviceDefSvc DeviceDefinitionsAPIService
	integrations []*grpc.Integration
	autoPiSvc    AutoPiAPIService
	deviceSvc    DeviceAPIService
	memoryCache  *gocache.Cache
}

func NewDeviceStatusIngestService(db func() *db.ReaderWriter, log *zerolog.Logger, eventService EventService, ddSvc DeviceDefinitionsAPIService, autoPiSvc AutoPiAPIService, deviceSvc DeviceAPIService) *DeviceStatusIngestService {
	// Cache the list of integrations.
	integrations, err := ddSvc.GetIntegrations(context.Background())
	if err != nil {
		log.Fatal().Err(err).Str("func", "NewDeviceStatusIngestService").Msg("Couldn't retrieve global integration list.")
	}
	c := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids

	return &DeviceStatusIngestService{
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
func (i *DeviceStatusIngestService) ProcessDeviceStatusMessages(ctx goka.Context, msg interface{}) {
	if err := i.processMessage(ctx, msg.(*DeviceStatusEvent)); err != nil {
		i.log.Err(err).Msg("Error processing device status message.")
	}
}

func (i *DeviceStatusIngestService) processMessage(ctx goka.Context, event *DeviceStatusEvent) error {
	if event.Type != deviceStatusEventType {
		return fmt.Errorf("received vehicle status event with unexpected type %s", event.Type)
	}

	integration, err := i.getIntegrationFromEvent(event)
	if err != nil {
		return err
	}

	switch integration.Vendor {
	case constants.SmartCarVendor:
		defer appmetrics.SmartcarIngestTotalOps.Inc()
	case constants.AutoPiVendor:
		defer appmetrics.AutoPiIngestTotalOps.Inc()
	}

	return i.processEvent(ctx, event)
}

var userDeviceDataPrimaryKeyColumns = []string{models.UserDeviceDatumColumns.UserDeviceID, models.UserDeviceDatumColumns.IntegrationID}

func cacheKey(userDeviceID, integrationID string) string {
	return fmt.Sprintf("%s_%s", userDeviceID, integrationID)
}

// processEvent handles the device data status update so we have a latest snapshot and saves to signals. This should all be refactored to device data api.
func (i *DeviceStatusIngestService) processEvent(_ goka.Context, event *DeviceStatusEvent) error {
	ctx := context.Background() // todo: will this still work with goka context instead?
	userDeviceID := event.Subject

	integration, err := i.getIntegrationFromEvent(event)
	if err != nil {
		return err
	}

	ck := cacheKey(userDeviceID, integration.Id)

	var device *pb.UserDevice
	maybeDevice, ok := i.memoryCache.Get(ck)
	if ok {
		device = maybeDevice.(*pb.UserDevice)
	} else {
		device, err = i.deviceSvc.GetUserDevice(ctx, userDeviceID)
		if err != nil {
			return fmt.Errorf("failed to find user_device with id %s. error: %w", userDeviceID, err)
		}

		// Validate integration Id
		i.memoryCache.Set(ck, device, 30*time.Minute)
	}

	var apiIntegration *pb.UserDeviceIntegration

	for _, connIntegration := range device.Integrations {
		if connIntegration.Id == integration.Id {
			apiIntegration = connIntegration
			break
		}
	}

	if apiIntegration == nil {
		return fmt.Errorf("device %s got a status for unpaired integration %s", userDeviceID, integration.Id)
	}

	deviceDefinitionResponse, err := i.deviceDefSvc.GetDeviceDefinitionByID(ctx, device.DeviceDefinitionId)
	if err != nil {
		return fmt.Errorf("failed to look up device definition %s: %w", device.DeviceDefinitionId, err)
	}

	if deviceDefinitionResponse == nil {
		return fmt.Errorf("failed to look up device definition %s", device.DeviceDefinitionId)
	}

	// update status to Active if not already set
	if apiIntegration.Status != constants.UserDeviceAPIIntegrationStatusActive {
		i.log.Info().Str("user_device_id", userDeviceID).Str("integration_id", apiIntegration.Id).
			Msgf("updating integration status for user device. current status %s", apiIntegration.Status)
		apiIntegration.Status = constants.UserDeviceAPIIntegrationStatusActive
		if _, err := i.deviceSvc.UpdateStatus(ctx, userDeviceID, apiIntegration.Id, apiIntegration.Status); err != nil {
			return fmt.Errorf("failed to update API integration: %w", err)
		}

		if integration.Vendor == constants.AutoPiVendor {
			err := i.autoPiSvc.UpdateState(apiIntegration.ExternalId, apiIntegration.Status)
			if err != nil {
				return fmt.Errorf("failed to update status when calling autopi api for deviceId: %s", apiIntegration.ExternalId)
			}
		}

		i.memoryCache.Delete(userDeviceID + "_" + integration.Id)
	}

	// techdebt: could likely get rid of this with tweak in app so that people just see that data came through - not specific to odometer
	// Null for most AutoPis.
	var newOdometer null.Float64
	if o, err := extractOdometer(event.Data); err == nil {
		newOdometer = null.Float64From(o)
	} else if integration.Vendor == constants.AutoPiVendor {
		// For AutoPis, for the purpose of odometer events we are pretending to always have
		// an odometer reading. Users became accustomed to seeing the associated events, even
		// though we mostly don't have odometer readings for AutoPis. For now, we fake it.

		// Update PLA-934:  Now that we are starting to receive real odometer values from
		//             		the AutoPi, we need the real odometer timestamp. To avoid alarming
		//					users as mentioned above, we resolved to create another column
		//					called "real_last_odometer_event_at" to store this value
		newOdometer = null.Float64From(0.0)
	}

	var newLastLocation null.Float64
	if o, err := extractLastLocation(event.Data); err == nil {
		newLastLocation = null.Float64From(o)
	}

	var newObd2 null.Float64
	if o, err := extractObd2(event.Data); err == nil {
		newObd2 = null.Float64From(o)
	}

	var datum *models.UserDeviceDatum

	deviceData, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(userDeviceID),
		models.UserDeviceDatumWhere.IntegrationID.EQ(integration.Id),
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GT(time.Now().Add(-14*24*time.Hour)),
	).All(ctx, i.db().Reader)

	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}

	if len(deviceData) > 0 {
		// Update the existing record.
		datum = deviceData[0]
	} else {
		// Insert a new record.
		datum = &models.UserDeviceDatum{UserDeviceID: userDeviceID, IntegrationID: integration.Id}
		i.memoryCache.Delete(userDeviceID + "_" + integration.Id)
	}

	i.processOdometer(datum, newOdometer, device, deviceDefinitionResponse, integration.Id)

	i.processLastLocation(datum, newLastLocation)

	i.processObd2(datum, newObd2)

	// Not every update has every signal. Merge the new into the old.
	compositeData := make(map[string]any)

	// This will preserve any mappings with keys present in datum.Data but not in
	// event.Data. If a key is present in both maps then the value from event.Data
	// takes precedence.
	//
	// For example, if in the database we have {A: 1, B: 2} and the new event has
	// {B: 4, C: 9} then the result should be {A: 1, B: 4, C: 9}.
	if err := json.Unmarshal(event.Data, &compositeData); err != nil {
		return err
	}

	datum.ErrorData = null.JSON{}

	// extract signals with timestamps and persist to signals
	existingSignalData := make(map[string]any)
	if err := datum.Signals.Unmarshal(&existingSignalData); err != nil {
		return err
	}
	// unmarshall only the event data
	eventData := make(map[string]any)
	err = json.Unmarshal(event.Data, &eventData)
	if err != nil {
		return errors.Wrap(err, "could not unmarshall event data")
	}
	newSignals, err := mergeSignals(existingSignalData, eventData, event.Time)
	if err != nil {
		return err
	}
	if err := datum.Signals.Marshal(newSignals); err != nil {
		return err
	}

	if err := datum.Upsert(ctx, i.db().Writer, true, userDeviceDataPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf("error upserting datum: %w", err)
	}

	switch integration.Vendor {
	case constants.SmartCarVendor:
		appmetrics.SmartcarIngestSuccessOps.Inc()
	case constants.AutoPiVendor:
		appmetrics.AutoPiIngestSuccessOps.Inc()
	}

	return nil
}

// processOdometer emits an odometer event and updates the last_odometer_event timestamp on the
// data record if the following conditions are met:
//   - there is no existing timestamp, or an hour has passed since that timestamp,
//   - the incoming status update has an odometer value, and
//   - the old status update lacks an odometer value, or has an odometer value that differs from
//     the new odometer reading
func (i *DeviceStatusIngestService) processOdometer(datum *models.UserDeviceDatum, newOdometer null.Float64, device *pb.UserDevice, dd *grpc.GetDeviceDefinitionItemResponse, integrationID string) {
	if !newOdometer.Valid {
		return
	}

	var oldOdometer null.Float64
	var oldOdometerTimestamp null.Time
	if datum.Signals.Valid {
		if o, err := extractOdometer(datum.Signals.JSON); err == nil {
			oldOdometer = null.Float64From(o)
		}
		if t, err := extractOdometerTime(datum.Signals.JSON); err == nil {
			oldOdometerTimestamp = null.TimeFrom(t)
		}
	}

	now := time.Now()
	odometerOffCooldown := !oldOdometerTimestamp.Valid || now.Sub(oldOdometerTimestamp.Time) >= odometerCooldown
	odometerChanged := !oldOdometer.Valid || newOdometer.Float64 > oldOdometer.Float64

	if odometerOffCooldown && odometerChanged {
		oldOdometerTimestamp = null.TimeFrom(now)
		i.emitOdometerEvent(device, dd, integrationID, newOdometer.Float64)
	}
}

func (i *DeviceStatusIngestService) processLastLocation(datum *models.UserDeviceDatum, newLastLocation null.Float64) {
	if !newLastLocation.Valid {
		return
	}
	var oldLastLocation null.Float64

	if datum.Signals.Valid {
		if o, err := extractLastLocation(datum.Signals.JSON); err == nil {
			oldLastLocation = null.Float64From(o)
		}
	}

	now := time.Now()

	locationChanged := !oldLastLocation.Valid || newLastLocation.Float64 > oldLastLocation.Float64

	if locationChanged {
		datum.LastLocationEventAt = null.TimeFrom(now)
	}
}

func (i *DeviceStatusIngestService) processObd2(datum *models.UserDeviceDatum, newObd2 null.Float64) {
	if !newObd2.Valid {
		return
	}

	var oldObd2 null.Float64

	if datum.Signals.Valid {
		if o, err := extractLastLocation(datum.Signals.JSON); err == nil {
			oldObd2 = null.Float64From(o)
		}
	}

	now := time.Now()

	obd2Changed := !oldObd2.Valid || newObd2.Float64 > oldObd2.Float64

	if obd2Changed {
		datum.LastOdb2EventAt = null.TimeFrom(now)
	}

}

func (i *DeviceStatusIngestService) emitOdometerEvent(device *pb.UserDevice, dd *grpc.GetDeviceDefinitionItemResponse, integrationID string, odometer float64) {
	event := &Event{
		Type:    "com.dimo.zone.device.odometer.update",
		Subject: device.Id,
		Source:  "dimo/integration/" + integrationID,
		Data: OdometerEvent{
			Timestamp: time.Now(),
			UserID:    device.UserId,
			Device: odometerEventDevice{
				ID:    device.Id,
				Make:  dd.Make.Name,
				Model: dd.Type.Model,
				Year:  int(dd.Type.Year),
			},
			Odometer: odometer,
		},
	}
	if err := i.eventService.Emit(event); err != nil {
		i.log.Err(err).Msgf("Failed to emit odometer event for device %s", device.Id)
	}
}

func extractOdometer(data []byte) (float64, error) {
	result := gjson.GetBytes(data, "odometer")
	if !result.Exists() {
		return 0, errors.New("data payload did not have an odometer reading")
	}
	return result.Float(), nil
}

func extractOdometerTime(data []byte) (time.Time, error) {
	result := gjson.GetBytes(data, "odometer.timestamp")
	if !result.Exists() {
		return time.Time{}, errors.New("data payload did not have an odometer timestamp")
	}
	return result.Time(), nil
}

func extractLastLocation(data []byte) (float64, error) {
	result := gjson.GetBytes(data, "location")
	if !result.Exists() {
		return 0, errors.New("data payload did not have a last_location reading")
	}
	return result.Float(), nil
}

func extractObd2(data []byte) (float64, error) {

	possibleSignals := []string{"odometer", "speed", "engineLoad", "coolantTemp"}

	var result gjson.Result

	for _, signal := range possibleSignals {
		result = gjson.GetBytes(data, signal)
		if result.Exists() {
			break
		}
	}

	if !result.Exists() {
		return 0, errors.New("data payload did not have an obd2 reading")
	}

	return result.Float(), nil
}

func mergeSignals(currentData map[string]interface{}, newData map[string]interface{}, t time.Time) (map[string]interface{}, error) {

	merged := make(map[string]interface{})
	for k, v := range currentData {
		merged[k] = v
	}
	// now iterate over new data and update any keys present in the new data with the events timestamp
	for k, v := range newData {
		merged[k] = map[string]interface{}{
			"timestamp": t.Format("2006-01-02T15:04:05Z"), // utc tz RFC3339
			"value":     v,
		}
	}
	return merged, nil
}

func (i *DeviceStatusIngestService) getIntegrationFromEvent(event *DeviceStatusEvent) (*grpc.Integration, error) {
	for _, integration := range i.integrations {
		if strings.HasSuffix(event.Source, integration.Id) {
			return integration, nil
		}
	}
	return nil, fmt.Errorf("no matching integration found in DB for event source: %s", event.Source)
}

type odometerEventDevice struct {
	ID    string `json:"id"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

type OdometerEvent struct {
	Timestamp time.Time           `json:"timestamp"`
	UserID    string              `json:"userId"`
	Device    odometerEventDevice `json:"device"`
	Odometer  float64             `json:"odometer"`
}

type DeviceStatusEvent struct {
	ID          string          `json:"id"`
	Source      string          `json:"source"`
	Specversion string          `json:"specversion"`
	Subject     string          `json:"subject"`
	Time        time.Time       `json:"time"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
}
