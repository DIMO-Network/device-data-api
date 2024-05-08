package controllers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared/privileges"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/controllers"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/DIMO-Network/device-data-api/internal/test"
)

const (
	mockDeviceDataIndexName = "mock-es-device-data"
)

type DeviceDataControllerV2Suite struct {
	suite.Suite
	ctx       context.Context
	mockCtrl  *gomock.Controller
	deviceSvc *mock_services.MockDeviceAPIService
	esMock    *MockEsInterface
	SUT       controllers.DeviceDataControllerV2
}

func (d *DeviceDataControllerV2Suite) SetupSuite() {
	d.ctx = context.Background()
	d.mockCtrl = gomock.NewController(d.T())
	d.deviceSvc = mock_services.NewMockDeviceAPIService(d.mockCtrl)

	d.esMock = NewMockEsInterface(d.mockCtrl)

	logger := zerolog.Nop()

	d.SUT = controllers.NewDeviceDataControllerV2(&config.Settings{DeviceDataIndexName: mockDeviceDataIndexName}, &logger, d.deviceSvc, d.esMock, nil, nil)
}

func TestDeviceDataControllerV2Suite(t *testing.T) {
	suite.Run(t, new(DeviceDataControllerV2Suite))
}

func (d *DeviceDataControllerV2Suite) TestGetDistanceDriven() {
	expectedDeviceID := ksuid.New().String()
	tests := []struct {
		name             string
		urlParameters    map[string]string
		expectedResponse int
		expectedDevice   *grpc.UserDevice
		customClaims     []privileges.Privilege
	}{
		{
			name: "Test happy path odometer success",
			urlParameters: map[string]string{
				"tokenID": "123",
			},
			expectedResponse: fiber.StatusOK,
			expectedDevice: &grpc.UserDevice{
				Id: expectedDeviceID,
			},
			customClaims: []privileges.Privilege{privileges.VehicleNonLocationData},
		},
		{
			name: "Test no device for token error",
			urlParameters: map[string]string{
				"tokenID": "123",
			},
			expectedResponse: fiber.StatusBadRequest,
			expectedDevice:   nil,
			customClaims:     []privileges.Privilege{privileges.VehicleNonLocationData},
		},
	}
	d.esMock.EXPECT().GetTotalDistanceDriven(gomock.Any(), expectedDeviceID).Return([]byte(`{"aggregations": {"max_odometer": {"value": 200},"min_odometer": {"value": 50}}}`), nil).AnyTimes()

	testUserID := "123123"
	app := fiber.New()

	for _, tc := range tests {
		app.Use(test.ClaimsInjectorTestHandler(tc.customClaims))
		app.Get("/v2/vehicles/:tokenID/analytics/total-distance", test.AuthInjectorTestHandler(testUserID), d.SUT.GetDistanceDriven)

		var errDevice error
		if tc.expectedDevice == nil {
			errDevice = errors.New("device not found")
		}
		d.deviceSvc.EXPECT().GetUserDeviceByTokenID(gomock.Any(), int64(123)).Return(tc.expectedDevice, errDevice)

		request := test.BuildRequest("GET", "/v2/vehicles/"+tc.urlParameters["tokenID"]+"/analytics/total-distance", "")
		response, err := app.Test(request)
		d.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		d.Require().NoError(err)

		if tc.expectedResponse == fiber.StatusOK && tc.expectedResponse == response.StatusCode {
			odmRes := struct {
				DistanceDriven float64 `json:"distanceDriven"`
			}{}

			err = json.Unmarshal(body, &odmRes)
			d.Require().NoError(err)

			d.Assert().Equal(float64(150), odmRes.DistanceDriven)
		}

		d.Require().Equal(tc.expectedResponse, response.StatusCode, "")
	}
}

func (d *DeviceDataControllerV2Suite) TestGetDailyDistance() {
	expectedDeviceID := ksuid.New().String()
	tests := []struct {
		name             string
		urlParameters    map[string]string
		expectedResponse int
		expectedDevice   *grpc.UserDevice
	}{
		{
			name: "Test happy path daily distance success",
			urlParameters: map[string]string{
				"tokenID": "123",
			},
			expectedResponse: fiber.StatusOK,
			expectedDevice: &grpc.UserDevice{
				Id: expectedDeviceID,
			},
		},
		{
			name: "Test no device for token error",
			urlParameters: map[string]string{
				"tokenID": "123",
			},
			expectedResponse: fiber.StatusBadRequest,
			expectedDevice:   nil,
		},
	}

	d.esMock.EXPECT().GetTotalDailyDistanceDriven(gomock.Any(), "America/Los_Angeles", expectedDeviceID).Return([]byte(`{"aggregations":{"days":{"buckets":[{"key_as_string": "1712300869000","min_odom":{"value": 200},"max_odom": {"value": 500}}]}}}`), nil).AnyTimes()

	testUserID := "123123"
	app := fiber.New()
	app.Get("/v2/vehicles/:tokenID/analytics/daily-distance", test.AuthInjectorTestHandler(testUserID), d.SUT.GetDailyDistance)

	for _, tc := range tests {
		var errDevice error
		if tc.expectedDevice == nil {
			errDevice = errors.New("device not found")
		}
		d.deviceSvc.EXPECT().GetUserDeviceByTokenID(gomock.Any(), int64(123)).Return(tc.expectedDevice, errDevice)

		request := test.BuildRequest("GET", fmt.Sprintf("/v2/vehicles/%s/analytics/daily-distance?time_zone=%s", tc.urlParameters["tokenID"], url.QueryEscape("America/Los_Angeles")), "")
		response, err := app.Test(request)
		d.Require().NoError(err)

		body, err := io.ReadAll(response.Body)
		d.Require().NoError(err)

		if tc.expectedResponse == fiber.StatusOK && tc.expectedResponse == response.StatusCode {
			type days struct {
				Date     string
				Distance float64
			}
			resp := struct {
				Days []days
			}{}

			err = json.Unmarshal(body, &resp)
			d.Require().NoError(err)

			d.Assert().Equal([]days{
				{
					Date:     "1712300869",
					Distance: float64(300),
				},
			}, resp.Days)
		}

		d.Require().Equal(tc.expectedResponse, response.StatusCode, "")
	}
}
