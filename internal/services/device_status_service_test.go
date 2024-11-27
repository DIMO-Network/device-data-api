package services

import (
	"context"
	"testing"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/constants"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/DIMO-Network/device-data-api/internal/test"
	"github.com/DIMO-Network/device-data-api/models"
	ddgrpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/DIMO-Network/shared/privileges"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/volatiletech/null/v8"
	"go.uber.org/mock/gomock"
)

func Test_sortBySignalTimestampDesc(t *testing.T) {
	udd := models.UserDeviceDatumSlice{
		&models.UserDeviceDatum{
			UserDeviceID: "123",
			Signals: null.JSONFrom([]byte(`{ "odometer": {
    "value": 88164.32,
    "timestamp": "2023-04-27T15:57:37Z"
  }}`)),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			IntegrationID: "123",
		},
		&models.UserDeviceDatum{
			UserDeviceID: "123",
			Signals: null.JSONFrom([]byte(`{ "odometer": {
    "value": 88174.32,
    "timestamp": "2023-04-27T16:57:37Z"
  }}`)),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			IntegrationID: "345",
		},
	}
	// validate setup is ok
	assert.Equal(t, 88164.32, gjson.GetBytes(udd[0].Signals.JSON, "odometer.value").Float())
	assert.Equal(t, 88174.32, gjson.GetBytes(udd[1].Signals.JSON, "odometer.value").Float())
	// sort and validate
	sortBySignalTimestampDesc(udd, "odometer")
	assert.Equal(t, 88174.32, gjson.GetBytes(udd[0].Signals.JSON, "odometer.value").Float())
	assert.Equal(t, "2023-04-27T16:57:37Z", gjson.GetBytes(udd[0].Signals.JSON, "odometer.timestamp").Time().Format(time.RFC3339))
	assert.Equal(t, 88164.32, gjson.GetBytes(udd[1].Signals.JSON, "odometer.value").Float())
}

func Test_sortBySignalValueDesc(t *testing.T) {
	udd := models.UserDeviceDatumSlice{
		&models.UserDeviceDatum{
			UserDeviceID: "123",
			Signals: null.JSONFrom([]byte(`{ "odometer": {
    "value": 88164.32,
    "timestamp": "2023-04-27T15:57:37Z"
  }}`)),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			IntegrationID: "123",
		},
		&models.UserDeviceDatum{
			UserDeviceID: "123",
			Signals: null.JSONFrom([]byte(`{ "odometer": {
    "value": 88174.32,
    "timestamp": "2023-04-27T16:57:37Z"
  }}`)),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			IntegrationID: "345",
		},
	}
	// validate setup is ok
	assert.Equal(t, 88164.32, gjson.GetBytes(udd[0].Signals.JSON, "odometer.value").Float())
	assert.Equal(t, 88174.32, gjson.GetBytes(udd[1].Signals.JSON, "odometer.value").Float())
	// sort and validate
	sortBySignalValueDesc(udd, "odometer")
	assert.Equal(t, 88174.32, gjson.GetBytes(udd[0].Signals.JSON, "odometer.value").Float())
	assert.Equal(t, 88164.32, gjson.GetBytes(udd[1].Signals.JSON, "odometer.value").Float())
}

func Test_calculateRange(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(mockCtrl)

	deviceStatusSvc := NewDeviceStatusService(deviceDefSvc)

	ddID := ksuid.New().String()
	styleID := ksuid.New().String()
	attrs := []*ddgrpc.DeviceTypeAttribute{
		{
			Name:  "fuel_tank_capacity_gal",
			Value: "15",
		},
		{
			Name:  "mpg",
			Value: "20",
		},
	}
	deviceDefSvc.EXPECT().GetDeviceDefinitionBySlug(gomock.Any(), ddID).Times(2).Return(&ddgrpc.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: ddID,
		Verified:           true,
		DeviceAttributes:   attrs,
	}, nil)

	rge, err := deviceStatusSvc.CalculateRange(ctx, ddID, &styleID, .7)
	require.NoError(t, err)
	require.NotNil(t, rge)
	assert.Equal(t, 337.9614, *rge)

	rge, err = deviceStatusSvc.CalculateRange(ctx, ddID, &styleID, 70)
	require.NoError(t, err)
	require.NotNil(t, rge)
	assert.Equal(t, 337.9614, *rge)
}

func Test_deviceStatusService_PrepareDeviceStatusInformation(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(mockCtrl)
	udID := ksuid.New().String()
	autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
	smartCarInt := test.BuildIntegrationGRPC(constants.SmartCarVendor, 0, 0)
	dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "Mach E", 2020, autoPiInteg)

	// because we have range data from SC, below is not needed since CalculateRange won't be called, which calls Get DD
	// deviceDefSvc.EXPECT().GetDeviceDefinitionBySlug(gomock.Any(), dd[0].DeviceDefinitionId).Times(1).Return(dd[0], nil)

	deviceStatusSvc := NewDeviceStatusService(deviceDefSvc)

	// SC data setup to  older
	smartCarData := models.UserDeviceDatum{
		UserDeviceID: udID,
		Signals: null.JSONFrom([]byte(`{"oil": {"value": 0.6859999895095825, "timestamp": "2023-04-27T14:57:37Z"}, 
				"range": {"value": 187.79, "timestamp": "2023-04-27T14:57:37Z"}, 
				"tires": {"value":{"backLeft": 244, "backRight": 280, "frontLeft": 244, "frontRight": 252}, "timestamp": "2023-04-27T14:57:37Z"}, 
				"charging": {"value":false, "timestamp": "2023-04-27T14:57:37Z"}, 
				"latitude": {"value":33.675048828125, "timestamp": "2023-04-27T14:57:37Z"}, 
				"odometer": {"value":195677.59375, "timestamp": "2023-04-27T14:57:37Z"}, 
				"longitude": {"value":-117.85894775390625, "timestamp": "2023-04-27T14:57:37Z"}
				}`)),
		CreatedAt:           time.Now().Add(time.Minute * -5),
		UpdatedAt:           time.Now().Add(time.Minute * -5),
		LastLocationEventAt: null.TimeFrom(time.Now().Add(time.Minute * -5)),
		LastOdb2EventAt:     null.TimeFrom(time.Now().Add(time.Minute * -5)),
		IntegrationID:       smartCarInt.Id,
	}
	autoPiData := models.UserDeviceDatum{
		UserDeviceID: udID,
		Signals: null.JSONFrom([]byte(`{
				"latitude": { "value": 33.75, "timestamp": "2023-04-27T15:57:37Z" },
				"ambientTemp": { "value": 19, "timestamp": "2023-05-07T13:02:19Z" },
				"longitude": { "value": -117.91, "timestamp": "2023-04-27T15:57:37Z" } }`)),
		CreatedAt:     time.Now().Add(time.Minute * -1),
		UpdatedAt:     time.Now().Add(time.Minute * -1),
		IntegrationID: autoPiInteg.Id,
	}
	slice := models.UserDeviceDatumSlice{}
	slice = append(slice, &smartCarData)
	slice = append(slice, &autoPiData)

	snapshot := deviceStatusSvc.PrepareDeviceStatusInformation(ctx, slice, dd[0].DeviceDefinitionId, nil, []privileges.Privilege{1, 3, 4})

	assert.Equal(t, 187.79, *snapshot.Range)
	assert.Equal(t, false, *snapshot.Charging)
	assert.Equal(t, 244.0, snapshot.TirePressure.BackLeft)
	assert.Equal(t, 195677.59375, *snapshot.Odometer)
	assert.Equal(t, 19.0, *snapshot.AmbientTemp)
	assert.Equal(t, 33.75, *snapshot.Latitude, "expected autopi latitude")
	assert.Equal(t, -117.91, *snapshot.Longitude, "expected autopi longitude")
	assert.Equal(t, "2023-04-27T15:57:37Z", snapshot.RecordUpdatedAt.Format(time.RFC3339))
}
