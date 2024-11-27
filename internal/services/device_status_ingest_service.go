package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/DIMO-Network/devices-api/pkg/grpc"

	"github.com/tidwall/gjson"

	"github.com/DIMO-Network/device-data-api/models"
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
	autoPiSvc    AutoPiAPIService
	deviceSvc    DeviceAPIService
	memoryCache  *gocache.Cache
}

func NewDeviceStatusIngestService(db func() *db.ReaderWriter, log *zerolog.Logger, eventService EventService, autoPiSvc AutoPiAPIService, deviceSvc DeviceAPIService) *DeviceStatusIngestService {
	c := gocache.New(30*time.Minute, 60*time.Minute) // band-aid on top of band-aids

	return &DeviceStatusIngestService{
		db:           db,
		log:          log,
		eventService: eventService,
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

	var newLastLocation null.Float64
	if o, err := extractLastLocation(event.Data); err == nil {
		newLastLocation = null.Float64From(o)
	}

	var datum *models.UserDeviceDatum

	deviceData, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(userDeviceID),
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
		datum = &models.UserDeviceDatum{UserDeviceID: userDeviceID, IntegrationID: event.Source}
		i.memoryCache.Delete(userDeviceID + "_" + event.Source)
	}

	i.processLastLocation(datum, newLastLocation)

	i.processObd2(datum)

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
	// if autopi, do not ingest VIN, remove vin from eventData. vin comes from fingerprint now for AP.

	newSignals, err := mergeSignals(existingSignalData, eventData, event.Time, event.Source)
	if err != nil {
		return err
	}

	if err := datum.Signals.Marshal(newSignals); err != nil {
		return err
	}

	if err := datum.Upsert(ctx, i.db().Writer, true, userDeviceDataPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf("error upserting datum: %w", err)
	}

	return nil
}

// processOdometer emits an odometer event and updates the last_odometer_event timestamp on the
// data record if the following conditions are met:
//   - there is no existing timestamp, or an hour has passed since that timestamp,
//   - the incoming status update has an odometer value, and
//   - the old status update lacks an odometer value, or has an odometer value that differs from
//     the new odometer reading
func (i *DeviceStatusIngestService) processOdometer(datum *models.UserDeviceDatum, newOdometer null.Float64, device *pb.UserDevice, integrationID string) {
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
		i.emitOdometerEvent(device, integrationID, newOdometer.Float64)
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

func (i *DeviceStatusIngestService) processObd2(datum *models.UserDeviceDatum) {

	var obd2Exists bool

	if datum.Signals.Valid {
		if o, err := checkObd2Exists(datum.Signals.JSON); err == nil {
			obd2Exists = o
		}
	}

	now := time.Now()

	if obd2Exists {
		datum.LastOdb2EventAt = null.TimeFrom(now)
	}

}

func (i *DeviceStatusIngestService) emitOdometerEvent(device *pb.UserDevice, integrationID string, odometer float64) {
	event := &Event{
		Type:    "com.dimo.zone.device.odometer.update",
		Subject: device.Id,
		Source:  "dimo/integration/" + integrationID,
		Data: OdometerEvent{
			Timestamp: time.Now(),
			UserID:    device.UserId,
			Device: odometerEventDevice{
				ID: device.Id,
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

func checkObd2Exists(data []byte) (bool, error) {

	possibleSignals := []string{"odometer", "speed", "engineLoad", "coolantTemp"}

	var result gjson.Result

	for _, signal := range possibleSignals {
		result = gjson.GetBytes(data, signal)
		if result.Exists() {
			break
		}
	}

	if !result.Exists() {
		return false, errors.New("data payload did not have an obd2 reading")
	}

	return true, nil
}

func mergeSignals(currentData map[string]interface{}, newData map[string]interface{}, t time.Time, source string) (map[string]interface{}, error) {

	merged := make(map[string]interface{})
	for k, v := range currentData {
		merged[k] = v
	}
	// now iterate over new data and update any keys present in the new data with the events timestamp
	for k, v := range newData {
		merged[k] = map[string]interface{}{
			"timestamp": t.Format("2006-01-02T15:04:05Z"), // utc tz RFC3339
			"value":     v,
			"source":    source,
		}
	}
	return merged, nil
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
