package controllers

import (
	"context"
	_ "embed"
	"testing"

	"github.com/tidwall/sjson"

	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	pb "github.com/DIMO-Network/device-definitions-api/pkg/grpc"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestDeviceDataController_addRangeIfNotExists(t *testing.T) {
	controller := gomock.NewController(t)
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(controller)
	ddID := ksuid.New().String()

	deviceDefSvc.EXPECT().GetDeviceDefinition(gomock.Any(), ddID).Times(1).Return(&pb.GetDeviceDefinitionItemResponse{
		DeviceDefinitionId: ddID,
		Name:               "test car",
		Type:               nil,
		Verified:           true,
		Make:               nil,
		DeviceStyles:       nil,
		DeviceAttributes: []*pb.DeviceTypeAttribute{
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

	deviceDefSvc.EXPECT().GetDeviceDefinition(gomock.Any(), ddID).Times(0)

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
