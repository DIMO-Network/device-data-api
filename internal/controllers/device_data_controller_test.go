package controllers

import (
	"context"
	_ "embed"
	"math/big"

	"github.com/DIMO-Network/shared/privileges"

	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/constants"
	response2 "github.com/DIMO-Network/device-data-api/internal/response"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/DIMO-Network/device-data-api/internal/test"
	"github.com/DIMO-Network/device-data-api/models"
	ddgrpc "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	dagrpc "github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	smartcar "github.com/smartcar/go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/mock/gomock"
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
	deviceStatusSvc := mock_services.NewMockDeviceStatusService(mockCtrl)

	testUserID := "123123"
	c := NewDeviceDataController(&config.Settings{Port: "3000"}, &logger, deviceSvc, nil, deviceDefSvc, deviceStatusSvc, pdb.DBS)
	app := fiber.New()
	app.Get("/user/devices/:userDeviceID/status", test.AuthInjectorTestHandler(testUserID), c.GetUserDeviceStatus)

	t.Run("GET - device status", func(t *testing.T) {
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
			LastLocationEventAt: null.TimeFrom(time.Now().Add(time.Minute * -5)),
			LastOdb2EventAt:     null.TimeFrom(time.Now().Add(time.Minute * -5)),
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

		notCharging := false
		deviceStatusSvc.EXPECT().PrepareDeviceStatusInformation(gomock.Any(), gomock.Any(), dd[0].DeviceDefinitionId, nil, []privileges.Privilege{1, 3, 4}).Times(1).
			Return(response2.DeviceSnapshot{
				Charging:     &notCharging,
				Odometer:     getPtrFloat(195677.59375),
				Latitude:     getPtrFloat(33.75),
				Longitude:    getPtrFloat(-117.91),
				Range:        getPtrFloat(187.79),
				TirePressure: &smartcar.TirePressure{BackLeft: 244.0},
			})

		request := test.BuildRequest("GET", "/user/devices/"+udID+"/status", "")
		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)

		if assert.Equal(t, fiber.StatusOK, response.StatusCode) == false {
			fmt.Println("response body: " + string(body))
		}

		snapshot := new(response2.DeviceSnapshot)
		err = json.Unmarshal(body, snapshot)
		assert.NoError(t, err)

		assert.Equal(t, 187.79, *snapshot.Range)
		assert.Equal(t, false, *snapshot.Charging)
		assert.Equal(t, 244.0, snapshot.TirePressure.BackLeft)
		assert.Equal(t, 195677.59375, *snapshot.Odometer)
		assert.Equal(t, 33.75, *snapshot.Latitude, "expected autopi latitude")
		assert.Equal(t, -117.91, *snapshot.Longitude, "expected autopi longitude")
		//assert.Equal(t, "2023-04-27T15:57:37Z", snapshot.RecordUpdatedAt.Format(time.RFC3339))

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestUserDevicesController_GetVehicleStatusRaw(t *testing.T) {
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
	deviceStatusSvc := mock_services.NewMockDeviceStatusService(mockCtrl)

	testUserID := "123123"
	c := NewDeviceDataController(&config.Settings{Port: "3000"}, &logger, deviceSvc, nil, deviceDefSvc, deviceStatusSvc, pdb.DBS)
	app := fiber.New()

	// Custom Claims
	app.Use(test.ClaimsInjectorTestHandler(
		[]privileges.Privilege{privileges.VehicleNonLocationData},
	))

	app.Get("/vehicle/:tokenId/status-raw", test.AuthInjectorTestHandler(testUserID), c.GetVehicleStatusRaw)

	t.Run("GET - device raw status", func(t *testing.T) {

		// arrange db, insert some user_devices
		autoPiInteg := test.BuildIntegrationGRPC(constants.AutoPiVendor, 10, 0)
		smartCarInt := test.BuildIntegrationGRPC(constants.SmartCarVendor, 0, 0)
		dd := test.BuildDeviceDefinitionGRPC(ksuid.New().String(), "Ford", "Mach E", 2020, autoPiInteg)
		tokenID := new(big.Int)
		tokenID.SetString("123456789012345678901234567890", 10)
		udID := ksuid.New().String()
		deviceSvc.EXPECT().GetUserDeviceByTokenID(gomock.Any(), gomock.Any()).Times(1).Return(&dagrpc.UserDevice{
			Id:                 udID,
			DeviceDefinitionId: dd[0].DeviceDefinitionId,
			UserId:             testUserID,
		}, nil)

		userDeviceData := models.UserDeviceDatum{
			UserDeviceID: udID,
			Signals: null.JSONFrom([]byte(`{"soc": {"value": 0.78, "timestamp": "2023-10-17T05:22:20Z"}, 
				"vin": {"value": "XP7YGCEJ9PB148921", "timestamp": "2023-10-17T05:22:20Z"}, 
				"range": {"value": 316.46140416, "timestamp": "2023-10-17T05:22:20Z"}, 
				"speed": {"value": 32.18688, "timestamp": "2023-10-17T05:22:20Z"}, 
				"charger": {"value": {"power": 3}, "timestamp": "2023-10-16T14:12:42Z"}, 
				"charging": {"value": false, "timestamp": "2023-10-17T05:22:20Z"}, 
				"latitude": {"value": 53.729016, "timestamp": "2023-10-17T05:22:20Z"}, 
				"odometer": {"value": 9914.409470477953, "timestamp": "2023-10-17T05:22:20Z"}, 
				"longitude": {"value": 9.990799, "timestamp": "2023-10-17T05:22:20Z"}, 
				"timestamp": {"value": "2023-10-17T05:22:20.453508151Z", "timestamp": "2023-10-17T05:22:20Z"}, 
				"vehicleId": {"value": "929850482922516", "timestamp": "2023-10-17T05:22:20Z"}, 
				"ambientTemp": {"value": 6.5, "timestamp": "2023-10-17T05:22:20Z"}, 
				"chargeLimit": {"value": 1, "timestamp": "2023-10-17T05:22:20Z"}}`)),
			CreatedAt:           time.Now().Add(time.Minute * -5),
			UpdatedAt:           time.Now().Add(time.Minute * -5),
			LastLocationEventAt: null.TimeFrom(time.Now().Add(time.Minute * -5)),
			LastOdb2EventAt:     null.TimeFrom(time.Now().Add(time.Minute * -5)),
			IntegrationID:       smartCarInt.Id,
		}
		err := userDeviceData.Insert(ctx, pdb.DBS().Writer, boil.Infer())
		assert.NoError(t, err)

		request := test.BuildRequest("GET", "/vehicle/"+tokenID.String()+"/status-raw", "")

		response, _ := app.Test(request)
		body, _ := io.ReadAll(response.Body)
		if assert.Equal(t, fiber.StatusOK, response.StatusCode) == false {
			fmt.Println("response body: " + string(body))
		}
		fmt.Println("response body: " + string(body))

		jsonString := string(body)

		// assert NonLocationData to exist but location data to be removed
		assert.True(t, gjson.Get(jsonString, "soc").Exists())
		assert.True(t, gjson.Get(jsonString, "charging").Exists())
		assert.False(t, gjson.Get(jsonString, "latitude").Exists())
		assert.False(t, gjson.Get(jsonString, "longitude").Exists())

		//teardown
		test.TruncateTables(pdb.DBS().Writer.DB, t)
	})
}

func TestParseDateRange(t *testing.T) {
	for _, c := range []struct {
		name          string
		originalStart string
		originalEnd   string
		expectedStart string
		expectedEnd   string
		valid         bool
	}{
		{
			name:          "dateOnly",
			originalStart: "2023-05-04",
			originalEnd:   "2023-05-06",
			valid:         true,
			expectedStart: "2023-05-04",
			expectedEnd:   "2023-05-06",
		},
		{
			name:          "rfc3339",
			originalStart: "2023-05-04T09:00:00Z",
			originalEnd:   "2023-05-06T23:00:00Z",
			valid:         true,
			expectedStart: "2023-05-04T09:00:00Z",
			expectedEnd:   "2023-05-06T23:00:00Z",
		},
		{
			name:          "noValuesPassed",
			originalStart: "",
			originalEnd:   "",
			valid:         true,
		},
		{
			originalStart: "2023-05-04T09:00:00",
			originalEnd:   "2023-05-06T23:00:00",
		},
		{
			originalStart: "2023-05-04T09:00:00Z",
			originalEnd:   "2023-05-06T23:00:00",
		},
	} {
		parsedStart, parsedEnd, err := parseDateRange(c.originalStart, c.originalEnd)
		if !c.valid {
			assert.Error(t, err)
		} else {
			if c.name == "noValuesPassed" {
				assert.True(t, validDate(parsedStart))
				assert.True(t, validDate(parsedEnd))
				continue
			}
			assert.NoError(t, err)
			assert.Equal(t, c.expectedStart, parsedStart)
			assert.Equal(t, c.expectedEnd, parsedEnd)
		}
	}
}

func getPtrFloat(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}
