package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/pkg/errors"
)

//go:generate mockgen -source autopi_api_service.go -destination mocks/autopi_api_service_mock.go
type AutoPiAPIService interface {
	UpdateState(deviceID string, state string) error
}

type autoPiAPIService struct {
	Settings   *config.Settings
	httpClient shared.HTTPClientWrapper
	dbs        func() *db.ReaderWriter
}

var ErrNotFound = errors.New("not found")

func NewAutoPiAPIService(settings *config.Settings, dbs func() *db.ReaderWriter) AutoPiAPIService {
	h := map[string]string{"Authorization": "APIToken " + settings.AutoPiAPIToken}
	hcw, _ := shared.NewHTTPClientWrapper(settings.AutoPiAPIURL, "", 10*time.Second, h, true) // ok to ignore err since only used for tor check

	return &autoPiAPIService{
		Settings:   settings,
		httpClient: hcw,
		dbs:        dbs,
	}
}

// UpdateState calls https://api.dimo.autopi.io/dongle/devices/{DEVICE_ID}/ Note that the deviceID is the autoPi one.
func (a *autoPiAPIService) UpdateState(deviceID string, state string) error {
	userMetaDataStateInfo := make(map[string]interface{})
	userMetaDataStateInfo["state"] = state

	userMetaDataInfo := make(map[string]interface{})
	userMetaDataInfo["user_metadata"] = userMetaDataStateInfo

	payload, _ := json.Marshal(userMetaDataInfo)

	res, err := a.httpClient.ExecuteRequest(fmt.Sprintf("/dongle/devices/%s/", deviceID), "PATCH", payload)
	if err != nil {
		return errors.Wrapf(err, "error calling autopi api to path device %s", deviceID)
	}
	defer res.Body.Close() // nolint

	return nil
}
