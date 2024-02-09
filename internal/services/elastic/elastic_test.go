package elastic_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services/elastic"
	"github.com/DIMO-Network/device-data-api/internal/test/elasticcontainer"
	"github.com/DIMO-Network/shared/privileges"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
)

const deviceIndex = "vss-device-data"

var (
	testStartTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	//go:embed static_vehicle_data_test.json
	staticVehicleData []byte
)

// TestGetHistory tests the GetHistory function.
func TestGetHistory(t *testing.T) {
	ctx := context.Background()

	service, cleanup := setupService(ctx, t)
	t.Cleanup(cleanup)

	// insert mapping into Elasticsearch using a typed client
	err := elasticcontainer.AddVSSMapping(ctx, service.ESClient(), deviceIndex)
	require.NoErrorf(t, err, "could not add mapping: %v", err)
	loadStaticVehicleData(t, service.ESClient())

	testCases := []struct {
		name           string
		params         elastic.GetHistoryParams
		setDefaults    bool
		expectError    bool
		expectedIDs    []string
		expectedFields []string
		excludedFields []string
	}{
		// test cases for bucketing
		{
			name: "number of buckets equal to the interval",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      10,
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},

		{
			name: "number of buckets half of the interval",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      5,
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "2", "4", "6", "8"},
		},

		{
			name: "number of buckets double the interval",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      20,
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
		{
			name: "1 bucket",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      1,
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0"},
		},
		{
			name: "buckets + 1 returned from elastic",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime.Add(time.Millisecond * -1),
				Buckets:      1,
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0"},
		},
		{
			name:        "buckets set to 0 defaults to 1000",
			setDefaults: true,
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      0,
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
		{
			name:        "buckets set negative number, defaults to 1000",
			setDefaults: true,
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      -1,
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},

		// test cases for time parameters
		{
			name: "end time before start time",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      10,
				EndTime:      testStartTime.Add(time.Millisecond * -10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{},
		},
		{
			name:        "end time before start time negative bucket",
			setDefaults: true,
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      -1,
				EndTime:      testStartTime.Add(time.Millisecond * -10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{},
		},
		{
			name: "time range is 0",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      10,
				EndTime:      testStartTime,
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0"},
		},
		{
			name:        "start time is zero causing default to 14 days from start",
			setDefaults: true,
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    time.Time{},
				Buckets:      int(time.Hour * 24 * 14 / time.Millisecond), // enough buckets to get all the data
				EndTime:      testStartTime.Add(time.Millisecond * 10),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
		{
			name:        "end time is zero causing default to now",
			setDefaults: true,
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      int(time.Since(testStartTime) / time.Millisecond), // enough buckets to get all the data
				EndTime:      time.Time{},
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},

		// test case PrivilegeIDs filtering
		{
			name: "PrivilegeIDs unset",
			params: elastic.GetHistoryParams{
				DeviceID:  "1",
				StartTime: testStartTime,
				Buckets:   1,
				EndTime:   testStartTime.Add(time.Millisecond * 1),
			},
			// no fields exist in the result everything is excluded
			expectedIDs:    nil,
			excludedFields: []string{"data.misc.cell", "data.vehicle.currentLocation", "data.vehicle.powertrain"},
		},
		{
			name: "PrivilegeIDs set to AllTimeLocation",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      1,
				EndTime:      testStartTime.Add(time.Millisecond * 1),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleAllTimeLocation},
			},
			// only the location fields are present in the result
			expectedIDs:    nil, // no test ids are present in the result because the fields are excluded
			expectedFields: []string{"data.vehicle.currentLocation", "data.misc.cell"},
			excludedFields: []string{"data.vehicle.powertrain"},
		},
		{
			name: "PrivilegeIDs set to NonLocationData",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      1,
				EndTime:      testStartTime.Add(time.Millisecond * 1),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData},
			},
			// everything except the location fields are present in the result
			expectedIDs:    []string{"0"},
			expectedFields: []string{"data.vehicle.powertrain"},
			excludedFields: []string{"data.misc.cell", "data.vehicle.currentLocation"},
		},
		{
			name: "PrivilegeIDs set to Commands",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      1,
				EndTime:      testStartTime.Add(time.Millisecond * 1),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleCommands},
			},
			// no fields exist in the result everything is excluded
			expectedIDs:    nil,
			excludedFields: []string{"data.misc.cell", "data.vehicle.currentLocation", "data.vehicle.powertrain"},
		},
		{
			name: "PrivilegeIDs set to AllTimeLocation and NonLocationData",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      1,
				EndTime:      testStartTime.Add(time.Millisecond * 1),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			// everything is present in the result
			expectedIDs:    []string{"0"},
			expectedFields: []string{"data.vehicle.powertrain", "data.vehicle.currentLocation", "data.misc.cell"},
		},
		{
			name: "PrivilegeIDs set to AllTimeLocation and Commands",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    testStartTime,
				Buckets:      1,
				EndTime:      testStartTime.Add(time.Millisecond * 1),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleCommands, privileges.VehicleAllTimeLocation},
			},
			// only the location fields are present in the result
			expectedIDs:    nil, // no test ids are present in the result because the fields are excluded
			expectedFields: []string{"data.vehicle.currentLocation", "data.misc.cell"},
			excludedFields: []string{"data.vehicle.powertrain"},
		},

		// misc edge cases
		{
			name: "invalid bucket number",
			params: elastic.GetHistoryParams{
				DeviceID:  "1",
				StartTime: testStartTime,
				Buckets:   -1,
				EndTime:   testStartTime.Add(time.Millisecond * 10),
			},
			expectError: true,
		},
		{
			name: "extremely large time range single buckets",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    time.Unix(0, 1),
				Buckets:      1,
				EndTime:      time.Date(9999999, 1, 1, 0, 0, 0, 0, time.UTC),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0"},
		},
		{
			name: "extremely large time range and buckets",
			params: elastic.GetHistoryParams{
				DeviceID:     "1",
				StartTime:    time.Unix(0, 1),
				Buckets:      math.MaxInt64,
				EndTime:      time.Date(9999999, 1, 1, 0, 0, 0, 0, time.UTC),
				PrivilegeIDs: []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation},
			},
			expectedIDs: []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		},
	}

	// For each test case, call the GetHistory function with the custom parameters and check if the result is equal to the expected result.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setDefaults {
				tc.params.SetDefaultHistoryParams()
			}
			result, err := service.GetHistory(context.Background(), tc.params)
			if tc.expectError {
				require.Errorf(t, err, "expected error but got nil")
				return
			}
			require.NoErrorf(t, err, "could not get history: %v", err)

			// check if each id in the result is in the expected ids
			ids := []string{}
			for _, hit := range result {
				// using gjson parse the testId field and add it to the ids slice
				id := gjson.GetBytes(hit, "testID")
				if id.Exists() {
					ids = append(ids, id.String())
				}
				// verify that the expected fields are present in the result
				for _, field := range tc.expectedFields {
					fieldData := gjson.GetBytes(hit, field)
					require.Truef(t, fieldData.Exists(), "expected field %q not found in result", field)
				}
				for _, field := range tc.excludedFields {
					fieldData := gjson.GetBytes(hit, field)
					require.Falsef(t, fieldData.Exists(), "excluded field %q found in result", field)
				}
			}
			require.Truef(t, slices.Equal(tc.expectedIDs, ids), "expected ids: %v, got ids: %v", tc.expectedIDs, ids)
		})
	}
}

// setupService creates a new elastic container and a new elastic service with the client and the logger.
func setupService(ctx context.Context, t *testing.T) (*elastic.Service, func()) {
	t.Helper()
	settings, cleanup, err := elasticcontainer.Create(ctx)
	require.NoErrorf(t, err, "could not create elastic container: %v", err)

	cfg := config.Settings{}
	cfg.ElasticSearchAnalyticsHost = settings.Address
	cfg.ElasticSearchAnalyticsUsername = "elastic"
	cfg.ElasticSearchAnalyticsPassword = settings.Password
	cfg.DeviceDataIndexNameV2 = deviceIndex

	// Create a new elastic service with the client and the logger.
	logger := zerolog.Nop()
	service, err := elastic.New(&cfg, &logger, settings.CACert)
	if !assert.NoErrorf(t, err, "could not create elastic service: %v", err) {
		cleanup()
	}
	return service, cleanup
}

// loadStaticVehicleData marshals staticVehicleData into a slice of byte slices.
// and then inserts each byte slice into the elastic index.
func loadStaticVehicleData(t *testing.T, client *es8.TypedClient) {
	t.Helper()

	// masrshal staticVehicleData into a slice of byte slices
	data := []map[string]any{}
	if err := json.Unmarshal(staticVehicleData, &data); err != nil {
		t.Fatalf("could not unmarshal static vehicle data: %v", err)
	}

	// loop through the data and add override the time field and subject for each object so we have a consistent state for the test.
	// The time field is incremented by 1 millisecond for each object.
	i := 0
	for _, obj := range data {
		obj["time"] = testStartTime.Add(time.Millisecond * time.Duration(i)).UnixMilli()
		obj["subject"] = "1"
		// add new unique id for test validation
		obj["testID"] = strconv.Itoa(i)
		i++
		d, err := json.Marshal(obj)
		if err != nil {
			t.Fatalf("could not marshal static vehicle data: %v", err)
		}
		_, err = client.Index(deviceIndex).Refresh(refresh.True).Raw(bytes.NewReader(d)).Do(context.Background())
		if err != nil {
			t.Fatalf("could not insert static vehicle data: %v", err)
		}
	}
}
