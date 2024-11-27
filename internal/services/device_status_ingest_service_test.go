package services

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"

	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	pb "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/tidwall/gjson"
	"go.uber.org/mock/gomock"

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
				Id:         apInt.Id,
				Status:     "Active",
				ExternalId: "",
			},
		},
		Vin:                &vin,
		DeviceDefinitionId: deviceDefinitionID,
		VinConfirmed:       true,
	}, nil)
	deviceDefSvc.EXPECT().GetDeviceDefinitionBySlug(gomock.Any(), deviceDefinitionID).Return(&ddgrpc.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: deviceDefinitionID,
		Name:               "Malibu",
		Verified:           true,
		Model:              "Malibu",
		Year:               2012,
		Make: &ddgrpc.DeviceMake{
			Id:       ksuid.New().String(),
			Name:     "Chevrolet",
			NameSlug: "chevrolet",
		},
	}, nil)

	// Only making use the last parameter.
	integs, _ := deviceDefSvc.GetIntegrations(ctx)
	integrationID := integs[0].Id

	ingest := NewDeviceStatusIngestService(pdb.DBS, &logger, mes, deviceDefSvc, autoPISvc, deviceSvc)
	// add an existing autopi datum
	dat1 := models.UserDeviceDatum{
		UserDeviceID:        userDeviceID,
		Signals:             null.JSONFrom([]byte(`{"signal_name_version_1": {"timestamp": "xx", "value": 23.4}}`)),
		LastLocationEventAt: null.TimeFrom(time.Now().Add(-10 * time.Hour)),
		LastOdb2EventAt:     null.TimeFrom(time.Now().Add(-10 * time.Hour)),
		IntegrationID:       integrationID,
	}
	err := dat1.Insert(ctx, pdb.DBS().Writer, boil.Infer())
	require.NoError(t, err)
	// add a smartcar datum to make sure can handle multiple
	dat2 := models.UserDeviceDatum{
		UserDeviceID:        userDeviceID,
		Signals:             null.JSONFrom([]byte(`{"signal_name_version_2": {"timestamp": "xx", "value": 23.4}}`)),
		LastLocationEventAt: null.TimeFrom(time.Now().Add(-10 * time.Hour)),
		LastOdb2EventAt:     null.TimeFrom(time.Now().Add(-10 * time.Hour)),
		IntegrationID:       ksuid.New().String(), // just any other integrationId
	}
	err = dat2.Insert(ctx, pdb.DBS().Writer, boil.Infer())
	require.NoError(t, err)

	input := &DeviceStatusEvent{
		Source:      "dimo/integration/" + integrationID,
		Specversion: "1.0",
		Subject:     userDeviceID,
		Type:        deviceStatusEventType,
		Time:        time.Now(),
		Data:        []byte(`{"odometer": 45.22, "signal_name_version_2": 12.3}`),
	}

	var ctxGk goka.Context
	err = ingest.processEvent(ctxGk, input)
	require.NoError(t, err)

	// get updated dat1 from db
	updatedData, err := models.FindUserDeviceDatum(ctx, pdb.DBS().Reader, userDeviceID, integrationID)
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

// User device data is getting a different row for all incoming integrations
func TestUserDeviceIntegrationsDifferent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//assert := assert.New(t)

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

	autopiInt := test.BuildIntegrationDefaultGRPC("AutoPi", 10, 10, true)
	smartCarInt := test.BuildIntegrationDefaultGRPC("SmartCar", 10, 10, true)

	deviceDefSvc.EXPECT().GetIntegrations(gomock.Any()).Times(3).Return([]*ddgrpc.Integration{autopiInt, smartCarInt}, nil)

	deviceSvc.EXPECT().GetUserDevice(gomock.Any(), userDeviceID).Times(2).Return(&pb.UserDevice{
		Id:     userDeviceID,
		UserId: ksuid.New().String(),
		Integrations: []*pb.UserDeviceIntegration{
			{
				Id:         autopiInt.Id,
				Status:     "Active",
				ExternalId: "",
			},
			{
				Id:         smartCarInt.Id,
				Status:     "Active",
				ExternalId: "",
			},
		},
		Vin:                &vin,
		DeviceDefinitionId: deviceDefinitionID,
		VinConfirmed:       true,
	}, nil)

	deviceDefSvc.EXPECT().GetDeviceDefinitionBySlug(gomock.Any(), deviceDefinitionID).Times(2).Return(&ddgrpc.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: deviceDefinitionID,
		Name:               "Malibu",
		Verified:           true,
		Model:              "Malibu",
		Year:               2012,
		Make: &ddgrpc.DeviceMake{
			Id:       ksuid.New().String(),
			Name:     "Chevrolet",
			NameSlug: "chevrolet",
		},
	}, nil)

	// get all integrations
	integs, _ := deviceDefSvc.GetIntegrations(ctx)

	// add an existing autopi datum

	currentTime := time.Now()

	for _, integration := range integs {
		ingest := NewDeviceStatusIngestService(pdb.DBS, &logger, mes, deviceDefSvc, autoPISvc, deviceSvc)
		input := &DeviceStatusEvent{
			Source:      "dimo/integration/" + integration.Id,
			Specversion: "1.0",
			Subject:     userDeviceID,
			Type:        deviceStatusEventType,
			Time:        currentTime,
			Data:        []byte(`{"odometer": 45.22, "signal_name_version_2": 12.3}`),
		}

		var ctxGk goka.Context
		err := ingest.processEvent(ctxGk, input)
		require.NoError(t, err)
	}

	// get updated dat1 from db
	updatedDataAutoPi, err := models.FindUserDeviceDatum(ctx, pdb.DBS().Reader, userDeviceID, autopiInt.Id)
	require.NoError(t, err)

	// validate signals were updated, or not updated, as expected
	// assume UTC tz
	assert.Equal(t, currentTime.Format("2006-01-02T15:04:05Z"), gjson.GetBytes(updatedDataAutoPi.Signals.JSON, "odometer.timestamp").Str, "odometer ts should be updated from latest event")
	assert.Equal(t, 45.22, gjson.GetBytes(updatedDataAutoPi.Signals.JSON, "odometer.value").Num, "odometer value should be updated from latest event")

	updatedDataSmartCar, err := models.FindUserDeviceDatum(ctx, pdb.DBS().Reader, userDeviceID, smartCarInt.Id)
	require.NoError(t, err)

	// validate signals were updated, or not updated, as expected
	// assume UTC tz
	assert.Equal(t, currentTime.Format("2006-01-02T15:04:05Z"), gjson.GetBytes(updatedDataSmartCar.Signals.JSON, "odometer.timestamp").Str, "odometer ts should be updated from latest event")
	assert.Equal(t, 45.22, gjson.GetBytes(updatedDataSmartCar.Signals.JSON, "odometer.value").Num, "odometer value should be updated from latest event")

}

func TestDeviceStatusIngestService_processEvent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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

	autopiInt := test.BuildIntegrationDefaultGRPC("AutoPi", 10, 10, true)
	deviceDefSvc.EXPECT().GetIntegrations(gomock.Any()).Times(1).Return([]*ddgrpc.Integration{autopiInt}, nil)

	mes := &testEventService{
		Buffer: make([]*Event, 0),
	}

	ingest := NewDeviceStatusIngestService(pdb.DBS, &logger, mes, deviceDefSvc, autoPISvc, deviceSvc)
	udID := ksuid.New().String()
	deviceDefinitionID := ksuid.New().String()
	vin := "4T3R6RFVXMU023395"
	newVin := "4T3R6RFVXMU0233XX"

	userDeviceData := test.SetupCreateUserDeviceData(t, udID, autopiInt.Id, vin, pdb)
	assert.NotNil(t, userDeviceData)

	deviceSvc.EXPECT().GetUserDevice(gomock.Any(), udID).Return(&pb.UserDevice{
		Id:     udID,
		UserId: ksuid.New().String(),
		Integrations: []*pb.UserDeviceIntegration{
			{
				Id:         autopiInt.Id,
				Status:     "Active",
				ExternalId: "",
			},
		},
		Vin:                &vin,
		DeviceDefinitionId: deviceDefinitionID,
		VinConfirmed:       true,
	}, nil)

	deviceDefSvc.EXPECT().GetDeviceDefinitionBySlug(gomock.Any(), deviceDefinitionID).Return(&ddgrpc.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: deviceDefinitionID,
		Name:               "Malibu",
		Verified:           true,
		Model:              "Malibu",
		Year:               2012,
		Make: &ddgrpc.DeviceMake{
			Id:       ksuid.New().String(),
			Name:     "Chevrolet",
			NameSlug: "chevrolet",
		},
	}, nil)

	var ctxGk goka.Context

	err := ingest.processEvent(ctxGk, &DeviceStatusEvent{
		ID:          ksuid.New().String(),
		Source:      "dimo/integration/" + autopiInt.Id,
		Specversion: "1.0.0",
		Subject:     udID,
		Time:        time.Now().UTC(),
		Type:        deviceStatusEventType,
		Data:        []byte(`{"vin": "` + newVin + `","odometer": 42431}`),
	})
	assert.NoError(t, err)

	// todo: query models.user device data, and verify that signals was filled in
	updatedDataAutoPi, err := models.FindUserDeviceDatum(ctx, pdb.DBS().Reader, userDeviceData.UserDeviceID, userDeviceData.IntegrationID)
	require.NoError(t, err)

	assert.Equal(t, "", gjson.GetBytes(updatedDataAutoPi.Signals.JSON, "vin.value").Str)
	assert.Equal(t, 42431.0, gjson.GetBytes(updatedDataAutoPi.Signals.JSON, "odometer.value").Num)
	assert.Equal(t, "dimo/integration/"+autopiInt.Id, gjson.GetBytes(updatedDataAutoPi.Signals.JSON, "odometer.source").Str)
}
