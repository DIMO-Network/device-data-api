package rpc

import (
	"context"
	mock_services "github.com/DIMO-Network/device-data-api/internal/services/mocks"
	"github.com/DIMO-Network/device-data-api/internal/test"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"testing"
)

func Test_userDeviceData_GetSignals(t *testing.T) {
	// start database
	const migrationsDirRelPath = "../../migrations"
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
	deviceDefSvc := mock_services.NewMockDeviceDefinitionsAPIService(mockCtrl)
	deviceStatusSvc := mock_services.NewMockDeviceStatusService(mockCtrl)
	// dont need other deps
	uddApi := NewUserDeviceData(pdb.DBS, &logger, deviceDefSvc, deviceStatusSvc)
	// seed db with 3 different date_ids
	scIntId := ksuid.New().String()
	apIntId := ksuid.New().String()
	reportRow := &models.ReportVehicleSignalsEventsTracking{
		DateID:             "20230714",
		IntegrationID:      apIntId,
		DeviceMakeID:       "",
		PropertyID:         "",
		Model:              "",
		Year:               2021,
		DeviceDefinitionID: "",
		DeviceMake:         "",
		Count:              0,
	}
	err := reportRow.Insert(ctx, pdb.DBS().Writer, boil.Infer())
	require.NoError(t, err)
	reportRow.DateID = "20230713"
	err = reportRow.Insert(ctx, pdb.DBS().Writer, boil.Infer())
	require.NoError(t, err)
	reportRow.DateID = "20230711"
	err = reportRow.Insert(ctx, pdb.DBS().Writer, boil.Infer())
	require.NoError(t, err)
	reportRow.DateID = "20230710"
	reportRow.IntegrationID = scIntId
	err = reportRow.Insert(ctx, pdb.DBS().Writer, boil.Infer())
	require.NoError(t, err)

	// call and verify
	dates, err := uddApi.GetAvailableDates(ctx, nil)
	require.NoError(t, err)
	require.Len(t, dates.DateIds, 4)
	assert.Equal(t, "20230714", dates.DateIds[0].DateId)
	assert.Equal(t, apIntId, dates.DateIds[0].IntegrationId)
	assert.Equal(t, "20230713", dates.DateIds[1].DateId)
	assert.Equal(t, "20230711", dates.DateIds[2].DateId)
	assert.Equal(t, "20230710", dates.DateIds[3].DateId)
	assert.Equal(t, scIntId, dates.DateIds[3].IntegrationId)
}
