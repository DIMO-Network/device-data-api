package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	pr "github.com/DIMO-Network/shared/middleware/privilegetoken"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/null/v8"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/response"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/device-data-api/models"
)

type DeviceDataControllerV2 struct {
	Settings        *config.Settings
	log             *zerolog.Logger
	deviceAPI       services.DeviceAPIService
	esService       EsInterface
	dbs             func() *db.ReaderWriter
	deviceStatusSvc services.DeviceStatusService
}

// NewDeviceDataControllerV2 constructor
func NewDeviceDataControllerV2(settings *config.Settings, logger *zerolog.Logger, deviceAPIService services.DeviceAPIService, esService EsInterface, deviceStatusSvc services.DeviceStatusService, dbs func() *db.ReaderWriter) DeviceDataControllerV2 {
	return DeviceDataControllerV2{
		Settings:        settings,
		log:             logger,
		deviceAPI:       deviceAPIService,
		esService:       esService,
		dbs:             dbs,
		deviceStatusSvc: deviceStatusSvc,
	}
}

// GetDailyDistance godoc
// @Description  Get kilometers driven for a tokenID each day.
// @Tags         device-data
// @Produce      json
// @Success      200 {object} controllers.DailyDistanceResp
// @Failure      404 "no device found for user with provided parameters"
// @Param        tokenID  path   int  true   "token id"
// @Param	     time_zone query string true "IANAS time zone id, e.g., America/Los_Angeles"
// @Security     BearerAuth
// @Router       /v2/vehicle/{tokenID}/analytics/daily-distance [get]
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
// @Description  Get kilometers driven for a tokenID since connected (ie. since we have data available)
// @Description  if it returns 0 for distanceDriven it means we have no odometer data.
// @Tags         device-data
// @Produce      json
// @Success      200
// @Failure      404 "no device found for user with provided parameters"
// @Param        tokenID  path   string  true   "token id"
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

// GetVehicleStatus godoc
// @Description Returns the latest status update for the vehicle with a given token id.
// @Tags        device-data
// @Param       tokenId path int true "token id"
// @Produce     json
// @Success     200 {object} response.Device
// @Failure     404
// @Security    BearerAuth
// @Router      /v2/vehicle/{tokenId}/status [get]
func (d *DeviceDataControllerV2) GetVehicleStatus(c *fiber.Ctx) error {
	claims := c.Locals("tokenClaims").(pr.CustomClaims)
	privileges := claims.PrivilegeIDs

	tokenIDStr := c.Params("tokenID")
	tokenID, err := strconv.ParseInt(tokenIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Couldn't parse token id %q.", tokenIDStr))
	}

	// tid := pgtypes.NewNullDecimal(new(decimal.Big).SetBigMantScale(ti, 0))
	userDeviceNFT, err := d.deviceAPI.GetUserDeviceByTokenID(c.Context(), tokenID)
	if err != nil {
		return err
	}

	if userDeviceNFT == nil {
		return fiber.NewError(fiber.StatusNotFound, "NFT not found.")
	}

	deviceData, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(userDeviceNFT.Id),
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GT(time.Now().Add(-90*24*time.Hour)),
	).All(c.Context(), d.dbs().Reader)
	if err != nil {
		return err
	}

	if len(deviceData) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No status updates yet.")
	}

	ds := d.deviceStatusSvc.PrepareDeviceStatusInformation(c.Context(), deviceData, userDeviceNFT.DeviceDefinitionId, userDeviceNFT.DeviceStyleId, privileges)

	var dsv2 response.Device
	dsv2.RecordCreatedAt = ds.RecordCreatedAt
	dsv2.RecordUpdatedAt = ds.RecordUpdatedAt
	dsv2.Status.PowerTrain.TractionBattery.Charging.IsCharging = null.BoolFromPtr(ds.Charging)
	dsv2.Status.PowerTrain.FuelSystem.Level = null.Float64FromPtr(ds.FuelPercentRemaining)

	if ds.BatteryCapacity != nil {
		dsv2.Status.PowerTrain.TractionBattery.GrossCapacity = null.Float64From(float64(*ds.BatteryCapacity))
	}

	if ds.OilLevel != nil {
		switch ol := *ds.OilLevel; {
		case ol > 0.75:
			dsv2.Status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("CRITICALLY_HIGH")
		case ol >= 0.5:
			dsv2.Status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("NORMAL")
		case ol > 0.25:
			dsv2.Status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("LOW_NORMAL")
		default:
			dsv2.Status.PowerTrain.CombustionEngine.EngineOilLevel = null.StringFrom("CRITICALLY_LOW")
		}
	}

	dsv2.Status.PowerTrain.TractionBattery.StateOfCharge.Displayed = multNull(null.Float64FromPtr(ds.StateOfCharge), 100)
	dsv2.Status.PowerTrain.TractionBattery.StateOfCharge.Current = multNull(null.Float64FromPtr(ds.StateOfCharge), 100)
	dsv2.Status.PowerTrain.TractionBattery.Charging.ChargeLimit = multNull(null.Float64FromPtr(ds.ChargeLimit), 100)

	dsv2.Status.TravelledDistance = null.Float64FromPtr(ds.Odometer)
	dsv2.Status.PowerTrain.Transmission.TravelledDistance = null.Float64FromPtr(ds.Odometer)
	dsv2.Status.PowerTrain.Range = null.Float64FromPtr(ds.Range)
	dsv2.Status.PowerTrain.FuelSystem.Range = null.Float64FromPtr(ds.Range)
	dsv2.Status.LowVoltageBattery.CurrentVoltage = null.Float64FromPtr(ds.BatteryVoltage)
	dsv2.Status.Exterior.AirTemperature = null.Float64FromPtr(ds.AmbientTemp)

	if ds.TirePressure != nil {
		dsv2.Status.Chassis.Axle.Row1.Wheel.Left.Tire.Pressure = null.Float64From(ds.TirePressure.FrontLeft)
		dsv2.Status.Chassis.Axle.Row1.Wheel.Right.Tire.Pressure = null.Float64From(ds.TirePressure.FrontRight)
		dsv2.Status.Chassis.Axle.Row2.Wheel.Left.Tire.Pressure = null.Float64From(ds.TirePressure.BackLeft)
		dsv2.Status.Chassis.Axle.Row2.Wheel.Right.Tire.Pressure = null.Float64From(ds.TirePressure.BackRight)
	}

	dsv2.Status.CurrentLocation.Timestamp = null.StringFrom(ds.RecordUpdatedAt.Format(time.RFC3339))
	dsv2.Status.CurrentLocation.Latitude = null.Float64FromPtr(ds.Latitude)
	dsv2.Status.CurrentLocation.Longitude = null.Float64FromPtr(ds.Longitude)

	return c.JSON(dsv2)
}

func (d *DeviceDataControllerV2) getDeviceFromTokenID(ctx context.Context, tokenID string) (*grpc.UserDevice, error) {
	tkID, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		return nil, errors.New("could not process the provided tokenId")
	}

	device, err := d.deviceAPI.GetUserDeviceByTokenID(ctx, tkID)
	if err != nil {
		return nil, errors.New("could find a device using provided tokenId")
	}

	return device, nil
}
