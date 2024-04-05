package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
)

type DeviceDataControllerV2 struct {
	Settings        *config.Settings
	log             *zerolog.Logger
	deviceAPI       services.DeviceAPIService
	esService       EsInterface
	definitionsAPI  services.DeviceDefinitionsAPIService
	deviceStatusSvc services.DeviceStatusService
	dbs             func() *db.ReaderWriter
}

// NewDeviceDataControllerV2 constructor
func NewDeviceDataControllerV2(settings *config.Settings, logger *zerolog.Logger, deviceAPIService services.DeviceAPIService, esService EsInterface) DeviceDataControllerV2 {
	return DeviceDataControllerV2{
		Settings:  settings,
		log:       logger,
		deviceAPI: deviceAPIService,
		esService: esService,
	}
}

// GetDailyDistance godoc
// @Description  Get kilometers driven for a userDeviceID each day.
// @Tags         device-data
// @Produce      json
// @Success      200 {object} controllers.DailyDistanceResp
// @Failure      404 "no device found for user with provided parameters"
// @Param        userDeviceID  path   string  true   "user device id"
// @Param	     time_zone query string true "IANAS time zone id, e.g., America/Los_Angeles"
// @Security     BearerAuth
// @Router       /v2/vehicles/{tokenID}/analytics/daily-distance [get]
func (d *DeviceDataControllerV2) GetDailyDistance(c *fiber.Ctx) error {
	tz := c.Query("time_zone")
	tkID := c.Params("tokenID")

	device, err := d.getDeviceFromTokenID(c.Context(), tkID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could find a device using provided tokenId!")
	}

	res, err := d.esService.GetTotalDailyDistanceDriven(c.Context(), tz, device.Id)
	if err != nil {
		return fmt.Errorf("failed to get total distance driven for device with token id %s - %w", tkID, err)
	}

	var ddr dailyDistanceElasticResult

	err = json.Unmarshal(res, &ddr)
	if err != nil {
		return fmt.Errorf("failed to parse response %w", err)
	}

	buckets := ddr.Aggregations.Days.Buckets

	days := make([]DailyDistanceDay, len(buckets))

	for i, b := range buckets {
		var dp *float64

		if b.MaxOdom.Value != nil {
			if shared.IsOdometerValid(*b.MaxOdom.Value) {
				d := *b.MaxOdom.Value - *b.MinOdom.Value
				dp = &d
			}
		}

		day := DailyDistanceDay{
			Date:     buckets[i].KeyAsString[:10],
			Distance: dp,
		}

		days[i] = day
	}

	return c.JSON(DailyDistanceResp{Days: days})
}

// GetDistanceDriven godoc
// @Description  Get kilometers driven for a userDeviceID since connected (ie. since we have data available)
// @Description  if it returns 0 for distanceDriven it means we have no odometer data.
// @Tags         device-data
// @Produce      json
// @Success      200
// @Failure      404 "no device found for user with provided parameters"
// @Param        userDeviceID  path   string  true   "user device id"
// @Security     BearerAuth
// @Router       /v2/vehicles/{tokenID}/analytics/total-distance [get]
func (d *DeviceDataControllerV2) GetDistanceDriven(c *fiber.Ctx) error {
	tkID := c.Params("tokenID")

	device, err := d.getDeviceFromTokenID(c.Context(), tkID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "could find a device using provided tokenId!")
	}

	res, err := d.esService.GetTotalDistanceDriven(c.Context(), device.Id)
	if err != nil {
		return errors.Wrap(err, "error querying odometer")
	}

	var result odometerQueryResult

	err = json.Unmarshal(res, &result)
	if err != nil {
		return fmt.Errorf("failed to get distance driven for device with token id %s - %w", tkID, err)
	}

	endOdometer := 0.0
	startOdometer := 0.0

	if result.Aggregations.MaxOdometer.Value != nil {
		endOdometer = *result.Aggregations.MaxOdometer.Value
	}

	if result.Aggregations.MinOdometer.Value != nil {
		startOdometer = *result.Aggregations.MinOdometer.Value
	}

	distanceDriven := endOdometer - startOdometer

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"distanceDriven": distanceDriven,
		"units":          "kilometers",
	})
}

func (d *DeviceDataControllerV2) getDeviceFromTokenID(ctx context.Context, tokenID string) (*grpc.UserDevice, error) {
	tkID, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		return nil, errors.New("could not process the provided tokenId!")
	}

	device, err := d.deviceAPI.GetUserDeviceByTokenID(ctx, tkID)
	if err != nil {
		return nil, errors.New("could find a device using provided tokenId!")
	}

	return device, nil
}
