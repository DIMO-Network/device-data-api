package controllers

import (
	"context"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeviceDataController_addRangeIfNotExists(t *testing.T) {
	controller := gomock.NewController(t)
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(controller)
	ddID := ksuid.New().String()

	bodyWithRange, err := addRangeIfNotExists(context.Background(), deviceDefSvc, []byte(body), ddID)
	require.NoError(t, err)

}

// todo embed
var body = ``
