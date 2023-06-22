package services

import (
	"context"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/golang/mock/gomock"
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
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	assert := assert.New(t)

	mes := &testEventService{
		Buffer: make([]*Event, 0),
	}
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(mockCtrl)
	autoPISvc := mock_services.NewMockAutoPiAPIService(mockCtrl)
	deviceSvc := mock_services.NewMockDeviceAPIService(mockCtrl)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()
	userDeviceID := ksuid.New().String()
	deviceDefinitionID := ksuid.New().String()
	vin := "4T3R6RFVXMU023395"
	apInt := test.BuildIntegrationDefaultGRPC("AutoPi", 10, 10, true)
	deviceDefSvc.EXPECT().GetIntegrations(gomock.Any()).Times(2).Return([]*ddgrpc.Integration{apInt}, nil)
	deviceSvc.EXPECT().GetUserDevice(gomock.Any(), userDeviceID).Times(1).Return(&pb.UserDevice{
		Id:     userDeviceID,
		UserId: ksuid.New().String(),
		Integrations: []*pb.UserDeviceIntegration{
			{
				Id:         "2RZIQAmcSNHt0X6OhqEDFE1Wj0X",
				Status:     "Active",
				ExternalId: "",
			},
		},
		Vin:                &vin,
		DeviceDefinitionId: deviceDefinitionID,
		VinConfirmed:       true,
	}, nil)
	deviceDefSvc.EXPECT().GetDeviceDefinitionsByIDs(gomock.Any(), []string{deviceDefinitionID}).Return([]*ddgrpc.GetDeviceDefinitionItemResponse{
		{
			DeviceDefinitionId: deviceDefinitionID,
			Name:               "Malibu",
			Verified:           true,
			Type: &ddgrpc.DeviceType{
				Type:      "Vehicle",
				Make:      "Chevrolet",
				Model:     "Malibu",
				Year:      2012,
				MakeSlug:  "chevrolet",
				ModelSlug: "malibu",
			},
			Make: &ddgrpc.DeviceMake{
				Id:       ksuid.New().String(),
				Name:     "Chevrolet",
				NameSlug: "chevrolet",
			},
		},
	}, nil)

	// Only making use the last parameter.
	integs, _ := deviceDefSvc.GetIntegrations(ctx)
	integrationID := integs[0].Id

	ingest := NewDeviceStatusIngestService(pdb.DBS, &logger, mes, deviceDefSvc, autoPISvc, deviceSvc)

	dat1 := models.UserDeviceDatum{
		UserDeviceID:        userDeviceID,
		Signals:             null.JSONFrom([]byte(`{"signal_name_version_1": {"timestamp": "xx", "value": 23.4}}`)),
		LastOdometerEventAt: null.TimeFrom(time.Now().Add(-10 * time.Second)),
		IntegrationID:       null.StringFrom(integrationID),
	}

	input := &DeviceStatusEvent{
		Source:      "dimo/integration/" + integrationID,
		Specversion: "1.0",
		Subject:     userDeviceID,
		Type:        deviceStatusEventType,
		Time:        time.Now(),
		Data:        []byte(`{"odometer": 45.22, "signal_name_version_2": 12.3}`),
	}

	var ctxGk goka.Context
	err := ingest.processEvent(ctxGk, input)
	require.NoError(t, err)

	// get updated dat1 from db
	updatedData, err := models.FindUserDeviceDatum(ctx, pdb.DBS().Reader, userDeviceID)
	require.NoError(t, err)

	// validate signals were updated, or not updated, as expected
	assert.Equal("xx", gjson.GetBytes(dat1.Signals.JSON, "signal_name_version_1.timestamp").Str, "signal 1 ts should not change and be present")
	assert.Equal(23.4, gjson.GetBytes(dat1.Signals.JSON, "signal_name_version_1.value").Num, "signal 1 value should not change and be present")
	// assume UTC tz
	assert.Equal(input.Time.Format("2006-01-02T15:04:05Z"), gjson.GetBytes(updatedData.Signals.JSON, "odometer.timestamp").Str, "odometer ts should be updated from latest event")
	assert.Equal(45.22, gjson.GetBytes(updatedData.Signals.JSON, "odometer.value").Num, "odometer value should be updated from latest event")

	assert.Equal(input.Time.Format("2006-01-02T15:04:05Z"), gjson.GetBytes(updatedData.Signals.JSON, "signal_name_version_2.timestamp").Str, "signal 2 ts should be updated from latest event")
	assert.Equal(12.3, gjson.GetBytes(updatedData.Signals.JSON, "signal_name_version_2.value").Num, "signal 2 value should be updated from latest event")
}
