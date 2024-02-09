package controllers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/controllers"
	"github.com/DIMO-Network/device-data-api/internal/services/elastic"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/DIMO-Network/device-data-api/internal/test"
	"github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared/privileges"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	historicalPath  = "/v2/vehicle/123/history"
	testTimeoutSecs = 15
)

var (
	testStartTime = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	testEndTime   = time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)
	testJSON      = []json.RawMessage{[]byte(`{"key": "value1"}`), []byte(`{"key": "value2"}`), []byte(`{"key": "value3"}`)}
)

func TestGetHistoricalPermissionedV2(t *testing.T) {
	tests := []struct {
		name             string
		queryParameters  map[string]string
		customClaims     []privileges.Privilege
		expectedReq      elastic.GetHistoryParams
		expectedResponse int
	}{
		{
			name:         "Test default startTime",
			customClaims: []privileges.Privilege{privileges.VehicleAllTimeLocation},
			queryParameters: map[string]string{
				"endTime": "2021-01-15T00:00:00Z",
				"buckets": "10",
			},
			expectedReq: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				EndTime:      testStartTime.Add(time.Hour * 24 * 14), // two weeks from startTime
				Buckets:      10,
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleAllTimeLocation},
			},
			expectedResponse: fiber.StatusOK,
		},
		{
			name: "Test default buckets",
			queryParameters: map[string]string{
				"startTime": "2021-01-01T00:00:00Z",
				"endTime":   "2021-01-02T00:00:00Z",
			},
			customClaims: []privileges.Privilege{privileges.VehicleAllTimeLocation},
			expectedReq: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				EndTime:      testEndTime,
				Buckets:      1000,
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleAllTimeLocation},
			},
			expectedResponse: fiber.StatusOK,
		},
		{
			name: "Test with startTime and endTime time and buckets",
			queryParameters: map[string]string{
				"startTime": "2021-01-01T00:00:00Z",
				"endTime":   "2021-01-02T00:00:00Z",
				"buckets":   "10",
			},
			customClaims: []privileges.Privilege{privileges.VehicleAllTimeLocation},
			expectedReq: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				EndTime:      testEndTime,
				Buckets:      10,
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleAllTimeLocation},
			},
			expectedResponse: fiber.StatusOK,
		},
		{
			name: "Test non-location data",
			queryParameters: map[string]string{
				"startTime": "2021-01-01T00:00:00Z",
				"endTime":   "2021-01-02T00:00:00Z",
			},
			customClaims: []privileges.Privilege{privileges.VehicleNonLocationData},
			expectedReq: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				EndTime:      testEndTime,
				Buckets:      1000,
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData},
			},
			expectedResponse: fiber.StatusOK,
		},

		// test cases for errors
		{
			name: "Test with no claims",
			queryParameters: map[string]string{
				"startTime": "2021-01-01T00:00:00Z",
				"endTime":   "2021-01-02T00:00:00Z",
				"buckets":   "10",
			},
			customClaims: []privileges.Privilege{},
			expectedReq: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				EndTime:      testEndTime,
				Buckets:      10,
				PrivilegeIDs: []privileges.Privilege{},
			},
			expectedResponse: fiber.StatusUnauthorized,
		},
		{
			name:             "Test with invalid startTime time",
			queryParameters:  map[string]string{"startTime": "invalid"},
			customClaims:     []privileges.Privilege{privileges.VehicleAllTimeLocation},
			expectedResponse: fiber.StatusBadRequest,
		},
		{
			name:             "Test with invalid endTime time",
			queryParameters:  map[string]string{"endTime": "invalid"},
			customClaims:     []privileges.Privilege{privileges.VehicleAllTimeLocation},
			expectedResponse: fiber.StatusBadRequest,
		},
		{
			name:             "Test with invalid buckets",
			queryParameters:  map[string]string{"buckets": "invalid"},
			customClaims:     []privileges.Privilege{privileges.VehicleAllTimeLocation},
			expectedResponse: fiber.StatusBadRequest,
		},
		{
			name:             "Test negative buckets",
			queryParameters:  map[string]string{"buckets": "-1"},
			customClaims:     []privileges.Privilege{privileges.VehicleAllTimeLocation},
			expectedResponse: fiber.StatusBadRequest,
		},
	}

	ctrl := gomock.NewController(t)
	esMock := NewMockEsInterface(ctrl)
	deviceSvc := mock_services.NewMockDeviceAPIService(ctrl)
	deviceSvc.EXPECT().GetUserDeviceByTokenID(gomock.Any(), int64(123)).Return(&grpc.UserDevice{Id: "1"}, nil).AnyTimes()
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "devices-api").
		Logger()

	deviceCtrl := controllers.NewDeviceDataController(&config.Settings{Port: "3000"}, &logger, deviceSvc, esMock, nil, nil, nil)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reqURL := createURL(tc.queryParameters)
			req, err := http.NewRequest("GET", reqURL, nil)
			require.NoErrorf(t, err, "error creating request: %v", err)

			app := fiber.New()
			app.Use(test.ClaimsInjectorTestHandler(tc.customClaims))
			app.Get("/v2/vehicle/:tokenID/history", deviceCtrl.GetHistoricalPermissionedV2)

			// if we do not expect errors, we need to set the expected calls
			if tc.expectedResponse == fiber.StatusOK {
				esMock.EXPECT().GetHistory(gomock.Any(), tc.expectedReq).Return(testJSON, nil)
			}

			resp, err := app.Test(req, testTimeoutSecs)
			require.NoErrorf(t, err, "error testing request: %v", err)

			body, err := io.ReadAll(resp.Body)
			require.NoErrorf(t, err, "error reading response body: %v", err)
			defer resp.Body.Close()

			require.NotNil(t, resp)
			// if the expected response is 200 and the actual response is not 200 print the response body
			if tc.expectedResponse == fiber.StatusOK && resp.StatusCode != fiber.StatusOK {
				t.Logf("response body: %s", body)
			}
			require.Equalf(t, tc.expectedResponse, resp.StatusCode, "expected response code %d, got %d", tc.expectedResponse, resp.StatusCode)
			if resp.StatusCode == fiber.StatusOK {
				expectedJSON := `{"statuses":[` + string(testJSON[0]) + "," + string(testJSON[1]) + "," + string(testJSON[2]) + `]}`
				require.JSONEq(t, expectedJSON, string(body))
			}
		})
	}
}

// createURL creates a URL with query parameters and the tokenID
func createURL(queryParameters map[string]string) string {
	reqURL := historicalPath
	if len(queryParameters) > 0 {
		reqURL += "?"
		for k, v := range queryParameters {
			reqURL += fmt.Sprintf("%s=%s&", k, v)
		}
		reqURL = reqURL[:len(reqURL)-1]
	}
	return reqURL
}
