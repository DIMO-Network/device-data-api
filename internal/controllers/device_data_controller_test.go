package controllers

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/constants"
	"github.com/DIMO-Network/device-data-api/internal/test"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"io"
	"os"
	"testing"
	"time"

	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	ddgrpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	dagrpc "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func TestDeviceDataController_addRangeIfNotExists(t *testing.T) {
	controller := gomock.NewController(t)
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(controller)
	ddID := ksuid.New().String()

	deviceDefSvc.EXPECT().GetDeviceDefinitionByID(gomock.Any(), ddID).Times(1).Return(&ddgrpc.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: ddID,
		Name:               "test car",
		Type:               nil,
		Verified:           true,
		Make:               nil,
		DeviceStyles:       nil,
		DeviceAttributes: []*ddgrpc.DeviceTypeAttribute{
			{
				Name:  "mpg",
				Value: "30",
			},
			{
				Name:  "fuel_tank_capacity_gal",
				Value: "16",
			},
		},
	}, nil)

	bodyWithRange, err := addRangeIfNotExists(context.Background(), deviceDefSvc, []byte(elasticDeviceData), ddID, nil)
	require.NoError(t, err)

	range1 := gjson.GetBytes(bodyWithRange, "hits.hits.0._source.data.range").Num //0.9
	assert.Equal(t, 695.23488, range1)                                            // kilometers

	range2 := gjson.GetBytes(bodyWithRange, "hits.hits.1._source.data.range").Num //0.8
	assert.Equal(t, 617.98656, range2)                                            // kilometers

	range3 := gjson.GetBytes(bodyWithRange, "hits.hits.2._source.data.range").Num //0.7
	assert.Equal(t, 540.73824, range3)

	rangeSkip := gjson.GetBytes(bodyWithRange, "hits.hits.6._source.data.range").Num //0.6
	assert.Equal(t, 463.48992, rangeSkip)                                            // kilometers
}

//go:embed historical_data_test.json
var elasticDeviceData string

func TestDeviceDataController_addRangeIfNotExists_NoChangeIfRangeExists(t *testing.T) {
	controller := gomock.NewController(t)
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(controller)
	ddID := ksuid.New().String()

	deviceDefSvc.EXPECT().GetDeviceDefinitionByID(gomock.Any(), ddID).Times(0)

	// if range exists anywhere in the body, do not add range anywhere
	bodySetRange, err2 := sjson.Set(elasticDeviceData, "hits.hits.0._source.data.range", 100.50)
	require.NoError(t, err2)

	bodyWithRange, err := addRangeIfNotExists(context.Background(), deviceDefSvc, []byte(bodySetRange), ddID, nil)
	require.NoError(t, err)

	range1 := gjson.GetBytes(bodyWithRange, "hits.hits.0._source.data.range").Num
	assert.Equal(t, 100.50, range1) // kilometers
	range2 := gjson.GetBytes(bodyWithRange, "hits.hits.1._source.data.range").Num
	assert.Equal(t, float64(0), range2, "expected no range property to exist here")
}

func Test_removeOdometerIfInvalid(t *testing.T) {

	body := removeOdometerIfInvalid([]byte(elasticDeviceData))

	// check that all bad odometers removed
	odo2 := gjson.GetBytes(body, "hits.hits.2._source.data.odometer")
	assert.False(t, odo2.Exists())

	odo3 := gjson.GetBytes(body, "hits.hits.3._source.data.odometer")
	assert.False(t, odo3.Exists())

	odo4 := gjson.GetBytes(body, "hits.hits.4._source.data.odometer")
	assert.False(t, odo4.Exists())
}

const migrationsDirRelPath = "../../migrations"

func TestUserDevicesController_GetUserDeviceStatus(t *testing.T) {
	// arrange global db and route setup
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Logger()

	ctx := context.Background()
	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	deviceSvc := mock_services.NewMockDeviceAPIService(mockCtrl)
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(mockCtrl)

	testUserID := "123123"
	c := NewDeviceDataController(&config.Settings{Port: "3000"}, &logger, deviceSvc, nil, deviceDefSvc, pdb.DBS)
	app := fiber.New()
	app.Get("/user/devices/:userDeviceID/status", test.AuthInjectorTestHandler(testUserID), c.GetUserDeviceStatus)

	t.Run("GET - device status merge autopi and smartcar", func(t *testing.T) {
		// arrange db, insert some user_devices
		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		smartCarInt := test.BuildIntegrationGRPC(constants.SmartCarVendor, 0, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "Mach E", 2020, autoPiInteg)
		udID := ksuid.New().String()
		deviceSvc.EXPECT().GetUserDevice(gomock.Any(), udID).Times(1).Return(&dagrpc.UserDevice{
			Id:                 udID,
			DeviceDefinitionId: dd[0].DeviceDefinitionId,
			UserId:             testUserID,
		}, nil)
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
			LastOdometerEventAt: null.TimeFrom(time.Now().Add(time.Minute * -5)),
			IntegrationID:       smartCarInt.Id,
		}
		err := smartCarData.Insert(ctx, pdb.DBS().Writer, boil.Infer())
		assert.NoError(t, err)
		// newer autopi data, expect to replace lat/long
		autoPiData := models.UserDeviceDatum{
			UserDeviceID: udID,
			Signals: null.JSONFrom([]byte(`{"latitude": { "value": 33.75, "timestamp": "2023-04-27T15:57:37Z" },
				"longitude": { "value": -117.91, "timestamp": "2023-04-27T15:57:37Z" } }`)),
			CreatedAt:     time.Now().Add(time.Minute * -1),
			UpdatedAt:     time.Now().Add(time.Minute * -1),
			IntegrationID: autoPiInteg.Id,
		}
		err = autoPiData.Insert(ctx, pdb.DBS().Writer, boil.Infer())
		assert.NoError(t, err)

		request := test.BuildRequest("GET", "/user/devices/"+udID+"/status", "")
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		if assert.Equal(t, fiber.StatusOK, response.StatusCode) == false {
			fmt.Println("response body: " + string(body))
		}

		snapshot := new(DeviceSnapshot)
		err = json.Unmarshal(body, snapshot)
		assert.NoError(t, err)

		assert.Equal(t, 187.79, *snapshot.Range)
		assert.Equal(t, false, *snapshot.Charging)
		assert.Equal(t, 244.0, snapshot.TirePressure.BackLeft)
		assert.Equal(t, 195677.59375, *snapshot.Odometer)
		assert.Equal(t, 33.75, *snapshot.Latitude, "expected autopi latitude")
		assert.Equal(t, -117.91, *snapshot.Longitude, "expected autopi longitude")
		assert.Equal(t, "2023-04-27T15:57:37Z", snapshot.RecordUpdatedAt.Format(time.RFC3339))

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
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

func TestUserDevicesController_calculateRange(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.Background()
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(mockCtrl)
	deviceSvc := mock_services.NewMockDeviceAPIService(mockCtrl)

	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Logger()

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
	deviceDefSvc.EXPECT().GetDeviceDefinitionByID(gomock.Any(), ddID).Times(1).Return(&ddgrpc.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: ddID,
		Verified:           true,
		DeviceAttributes:   attrs,
	}, nil)

	pdb, container := test.StartContainerDatabase(ctx, t, migrationsDirRelPath)
	defer func() {
		ctx := context.Background()
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()
	_ = NewDeviceDataController(&config.Settings{Port: "3000"}, &logger, deviceSvc, nil, deviceDefSvc, pdb.DBS)
	rge, err := calculateRange(ctx, deviceDefSvc, ddID, &styleID, .7)
	require.NoError(t, err)
	require.NotNil(t, rge)
	assert.Equal(t, 337.9614, *rge)
}
