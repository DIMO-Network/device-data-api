package services

import (
	"context"
	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"os"
	"testing"
	"time"

	"github.com/tidwall/gjson"

	"github.com/DIMO-Network/device-data-api/internal/test"
	"github.com/DIMO-Network/device-data-api/models"
	ddgrpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/lovoo/goka"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
)

type testEventService struct {
	Buffer []*Event
}

func (e *testEventService) Emit(event *Event) error {
	e.Buffer = append(e.Buffer, event)
	return nil
}

const migrationsDirRelPath = "../../migrations"

// TestAutoPiStatus tests that the signals column is getting updated correctly merging any existing data and setting timestamps
func TestAutoPiStatus(t *testing.T) {
	assert := assert.New(t)

	mes := &testEventService{
		Buffer: make([]*Event, 0),
	}
	deviceDefSvc := testDeviceDefinitionSvc{}
	autoPISvc := testAutoPISvc{}
	deviceSvc := testDeviceSvc{}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	// Only making use the last parameter.
	ddID := ksuid.New().String()
	integs, _ := deviceDefSvc.GetIntegrations(ctx)
	integrationID := integs[0].Id

	ingest := NewDeviceStatusIngestService(pdb.DBS, &logger, mes, deviceDefSvc, autoPISvc, deviceSvc)

	dat1 := models.UserDeviceDatum{
		UserDeviceID:        ddID,
		Signals:             null.JSONFrom([]byte(`{"signal_name_version_1": {"timestamp": "xx", "value": 23.4}}`)),
		LastOdometerEventAt: null.TimeFrom(time.Now().Add(-10 * time.Second)),
		IntegrationID:       null.StringFrom(integrationID),
	}

	input := &DeviceStatusEvent{
		Source:      "dimo/integration/" + integrationID,
		Specversion: "1.0",
		Subject:     ddID,
		Type:        deviceStatusEventType,
		Time:        time.Now(),
		Data:        []byte(`{"odometer": 45.22, "signal_name_version_2": 12.3}`),
	}

	var ctxGk goka.Context
	err := ingest.processEvent(ctxGk, input)
	require.NoError(t, err)

	// validate signals were updated, or not updated, as expected
	assert.Equal("xx", gjson.GetBytes(dat1.Signals.JSON, "signal_name_version_1.timestamp").Str, "signal 1 ts should not change and be present")
	assert.Equal(23.4, gjson.GetBytes(dat1.Signals.JSON, "signal_name_version_1.value").Num, "signal 1 value should not change and be present")
	// assume UTC tz
	assert.Equal(input.Time.Format("2006-01-02T15:04:05Z"), gjson.GetBytes(dat1.Signals.JSON, "odometer.timestamp").Str, "odometer ts should be updated from latest event")
	assert.Equal(45.22, gjson.GetBytes(dat1.Signals.JSON, "odometer.value").Num, "odometer value should be updated from latest event")

	assert.Equal(input.Time.Format("2006-01-02T15:04:05Z"), gjson.GetBytes(dat1.Signals.JSON, "signal_name_version_2.timestamp").Str, "signal 2 ts should be updated from latest event")
	assert.Equal(12.3, gjson.GetBytes(dat1.Signals.JSON, "signal_name_version_2.value").Num, "signal 2 value should be updated from latest event")
}

type testAutoPISvc struct {
}

// nolint
func (t testAutoPISvc) UpdateState(deviceID string, state string) error {
	//TODO implement me
	return nil
}

type testDeviceDefinitionSvc struct {
}

// nolint
func (t testDeviceDefinitionSvc) GetDeviceDefinition(ctx context.Context, id string) (*ddgrpc.GetDeviceDefinitionItemResponse, error) {
	//TODO implement me
	return &ddgrpc.GetDeviceDefinitionItemResponse{}, nil
}

// nolint
func (t testDeviceDefinitionSvc) GetIntegrations(ctx context.Context) ([]*ddgrpc.Integration, error) {
	//TODO implement me
	return nil, nil
}

// nolint
func (t testDeviceDefinitionSvc) GetDeviceDefinitionsByIDs(ctx context.Context, ids []string) ([]*ddgrpc.GetDeviceDefinitionItemResponse, error) {
	//TODO implement me
	return nil, nil
}

type testDeviceSvc struct {
}

// nolint
func (t testDeviceSvc) ListUserDevicesForUser(ctx context.Context, userID string) (*pb.ListUserDevicesForUserResponse, error) {
	//TODO implement me
	return nil, nil
}

// nolint
func (t testDeviceSvc) GetUserDevice(ctx context.Context, userDeviceID string) (*pb.UserDevice, error) {
	//TODO implement me
	return nil, nil
}

// nolint
func (t testDeviceSvc) UserDeviceBelongsToUserID(ctx context.Context, userID, userDeviceID string) (bool, error) {
	//TODO implement me
	return true, nil
}

// nolint
func (t testDeviceSvc) GetUserDeviceByTokenID(ctx context.Context, tokenID int64) (*pb.UserDevice, error) {
	//TODO implement me
	return nil, nil
}

// nolint
func (t testDeviceSvc) UpdateStatus(ctx context.Context, userDeviceID string, integrationID string, status string) (*pb.UserDevice, error) {
	//TODO implement me
	return nil, nil
}
