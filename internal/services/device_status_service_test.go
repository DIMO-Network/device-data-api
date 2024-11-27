package services

import (
	"testing"
	"time"

	"github.com/DIMO-Network/device-data-api/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/volatiletech/null/v8"
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
