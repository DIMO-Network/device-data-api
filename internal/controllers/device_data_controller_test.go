package controllers

import (
	"context"
	_ "embed"
	"testing"

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

	range1 := gjson.GetBytes(bodyWithRange, "hits.hits.1._source.data.range").Num
	//0.6745098039215687 fpr, or make it a simpler number to be easier
	assert.Equal(t, 617.98656, range1) // kilometers

	// todo - do math by hand make sure all good

	// todo - test the skip node, to validate we can set records further down when skipping an index

	// then another test for not modifying range if it is present on the second row
	// needs to be with different data
	//rangeUnchanged := gjson.GetBytes(bodyWithRange, "hits.hits.2._source.data.range").Num
	//assert.Equal(t, 400, rangeUnchanged)
}

//go:embed historical_data_test.json
var elasticDeviceData string
