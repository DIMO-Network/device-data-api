//go:generate mockgen -package controllers_test -destination es_mock_test.go github.com/DIMO-Network/device-data-api/internal/controllers EsInterface
package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/DIMO-Network/devices-api/pkg/grpc"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	pr "github.com/DIMO-Network/shared/middleware/privilegetoken"
	"github.com/DIMO-Network/shared/privileges"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/volatiletech/null/v8"
	"golang.org/x/exp/slices"

	"github.com/DIMO-Network/device-data-api/internal/config"
	_ "github.com/DIMO-Network/device-data-api/internal/response" // also needed for swagger gen
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/device-data-api/models"
)

type DeviceDataController struct {
	Settings        *config.Settings
	log             *zerolog.Logger
	deviceAPI       services.DeviceAPIService
	esService       EsInterface
	deviceStatusSvc services.DeviceStatusService
	dbs             func() *db.ReaderWriter
}

// EsInterface is an interface for the elastic service.
type EsInterface interface {
	GetTotalDistanceDriven(ctx context.Context, deviceID string) ([]byte, error)
	GetTotalDailyDistanceDriven(ctx context.Context, tz, deviceID string) ([]byte, error)
	ESClient() *elasticsearch.TypedClient
}

// NewDeviceDataController constructor
func NewDeviceDataController(settings *config.Settings, logger *zerolog.Logger, deviceAPIService services.DeviceAPIService, esService EsInterface, deviceStatusSvc services.DeviceStatusService, dbs func() *db.ReaderWriter) DeviceDataController {
	return DeviceDataController{
		Settings:        settings,
		log:             logger,
		deviceAPI:       deviceAPIService,
		esService:       esService,
		dbs:             dbs,
		deviceStatusSvc: deviceStatusSvc,
	}
}

// GetHistoricalRaw godoc
// @Description  Get all historical data for a userDeviceID, within start and end range
// @Tags         device-data
// @Produce      json
// @Success      200
// @Param        userDeviceID  path   string  true   "user id"
// @Param        startDate     query  string  false  "startDate eg 2022-01-02. if empty two weeks back"
// @Param        endDate       query  string  false  "endDate eg 2022-03-01. if empty today"
// @Security     BearerAuth
// @Router       /v1/user/device-data/{userDeviceID}/historical [get]
func (d *DeviceDataController) GetHistoricalRaw(c *fiber.Ctx) error {
	const dateLayout = "2006-01-02" // date layout support by elastic
	userDeviceID := c.Params("userDeviceID")
	startDate := c.Query("startDate")
	if startDate == "" {
		startDate = time.Now().Add(-1 * (time.Hour * 24 * 14)).Format(dateLayout) // if no start date default to 2 weeks
	} else {
		_, err := time.Parse(dateLayout, startDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}
	endDate := c.Query("endDate")
	if endDate == "" {
		endDate = time.Now().Format(dateLayout)
	} else {
		_, err := time.Parse(dateLayout, endDate)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	userDevice, err := d.deviceAPI.GetUserDevice(c.Context(), userDeviceID)
	if err != nil {
		return err
	}

	return d.getHistoryV1(c, userDevice, startDate, endDate, types.SourceFilter{})
}

// GetHistoricalRawPermissioned godoc
// @Description  Get all historical data for a tokenID, within start and end range
// @Tags         device-data
// @Produce      json
// @Success      200
// @Param        tokenID  path   int64  true   "token id"
// @Param        startDate     query  string  false  "startDate ex: 2022-01-02; or,  2022-01-02T09:00:00Z; if empty two weeks back"
// @Param        endDate       query  string  false  "endDate ex: 2022-03-01; or, 2023-03-01T09:00:00Z; if empty today"
// @Security     BearerAuth
// @Router       /v1/vehicle/{tokenID}/history [get]
func (d *DeviceDataController) GetHistoricalRawPermissioned(c *fiber.Ctx) error {
	tokenID := c.Params("tokenID")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	startDate, endDate, err := parseDateRange(startDate, endDate)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	i, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userDevice, err := d.deviceAPI.GetUserDeviceByTokenID(c.Context(), i)
	if err != nil {
		return err
	}

	privs := getPrivileges(c)

	var filter types.SourceFilter

	if slices.Contains(privs, privileges.VehicleAllTimeLocation) {
		filter.Includes = append(filter.Includes, "data.latitude", "data.longitude", "location", "data.cell", "cell", "data.timestamp")
	} else {
		filter.Excludes = append(filter.Excludes, "data.latitude", "data.longitude", "location", "data.cell", "cell")
	}

	if slices.Contains(privs, privileges.VehicleNonLocationData) {
		// Overrides the more limited Includes entries from above if the token also
		// has AllTimeLocation.
		filter.Includes = append(filter.Includes, "*")
	}
	return d.getHistoryV1(c, userDevice, startDate, endDate, filter)
}

func (d *DeviceDataController) getHistoryV1(c *fiber.Ctx, userDevice *grpc.UserDevice, startDate, endDate string, filter types.SourceFilter) error {
	var source types.SourceConfig = filter
	msm := types.MinimumShouldMatch(1)
	req := search.Request{
		Query: &types.Query{
			FunctionScore: &types.FunctionScoreQuery{
				Query: &types.Query{
					Bool: &types.BoolQuery{
						Filter: []types.Query{
							{Term: map[string]types.TermQuery{"subject": {Value: userDevice.Id}}},
							{Range: map[string]types.RangeQuery{"data.timestamp": types.DateRangeQuery{Gte: some.String(startDate), Lte: some.String(endDate)}}},
						},
						Should: []types.Query{
							{Exists: &types.ExistsQuery{Field: "data.odometer"}},
							{Exists: &types.ExistsQuery{Field: "data.latitude"}},
						},
						MinimumShouldMatch: &msm,
					},
				},
				Functions: []types.FunctionScore{{RandomScore: &types.RandomScoreFunction{}}},
			},
		},
		Size:    some.Int(1000),
		Source_: &source,
	}

	res, err := d.esService.ESClient().Search().Index(d.Settings.DeviceDataIndexName).Request(&req).Perform(c.Context())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	localLog := d.log.With().Str("userDeviceId", userDevice.Id).Str("deviceDefinitionId", userDevice.DeviceDefinitionId).Logger()
	if res.StatusCode >= fiber.StatusBadRequest {
		localLog.Error().Str("userDeviceId", userDevice.Id).Interface("response", res).Msgf("Got status code %d from Elastic.", res.StatusCode)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		localLog.Err(err).Str("userDeviceId", userDevice.Id).Msg("Failed to read Elastic response body.")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}

	b := len(body)
	t := time.Now()
	body = removeOdometerIfInvalid(body)
	if d := time.Since(t); d > time.Second {
		localLog.Info().Dur("duration", d).Msgf("Performed JSON operations on %d bytes of history.", b)
	}

	c.Set("Content-Type", fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
}

// removeOdometerIfInvalid removes the odometer json properties we consider invalid
func removeOdometerIfInvalid(body []byte) []byte {
	if len(body) == 0 {
		return body
	}

	resultData := gjson.GetBytes(body, "hits.hits.#._source.data")
	for i, r := range resultData.Array() {
		// note range is reported in km
		odoResult := r.Get("odometer")
		if odoResult.Exists() {
			if !shared.IsOdometerValid(odoResult.Float()) {
				// set json to remove?
				body, _ = sjson.DeleteBytes(body, fmt.Sprintf("hits.hits.%d._source.data.odometer", i))
			}
		}
	}

	return body
}

type odometerQueryResult struct {
	Aggregations struct {
		MaxOdometer struct {
			Value *float64 `json:"value"`
		} `json:"max_odometer"`
		MinOdometer struct {
			Value *float64 `json:"value"`
		} `json:"min_odometer"`
	} `json:"aggregations"`
}

// GetDistanceDriven godoc
// @Description  Get kilometers driven for a userDeviceID since connected (ie. since we have data available)
// @Description  [🔴__Warning - API Shutdown by June 30, 2024, Use `/v2/vehicles/:tokenId/analytics/total-distance` instead__🔴]  if it returns 0 for distanceDriven it means we have no odometer data.
// @Tags         device-data [End Of Life Warning]
// @Produce      json
// @Success      200
// @Failure      404 "no device found for user with provided parameters"
// @Param        userDeviceID  path   string  true   "user device id"
// @Security     BearerAuth
// @Router       /v1/user/device-data/{userDeviceID}/distance-driven [get]
func (d *DeviceDataController) GetDistanceDriven(c *fiber.Ctx) error {
	userDeviceID := c.Params("userDeviceID")
	query := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Term: map[string]types.TermQuery{"subject": {Value: userDeviceID}}},
				},
			},
		},
		Size: some.Int(0),
		Aggregations: map[string]types.Aggregations{
			"max_odometer": {
				Max: &types.MaxAggregation{
					Field: some.String("data.odometer"),
				},
			},
			"min_odometer": {
				Min: &types.MinAggregation{
					Field: some.String("data.odometer"),
				},
			},
		},
	}

	res, err := d.esService.ESClient().Search().Index(d.Settings.DeviceDataIndexName).Request(query).Perform(c.Context())
	if err != nil {
		return errors.Wrap(err, "error querying odometer")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status code %d from Elastic", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var result odometerQueryResult

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
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

// GetUserDeviceStatus godoc
// @Description Returns the latest status update for the device. May return 404 if the
// @Description user does not have a device with the ID, or if no status updates have come. Note this endpoint also exists under nft_controllers
// @Tags        device-data
// @Produce     json
// @Param       user_device_id path     string true "user device ID"
// @Success     200            {object} response.DeviceSnapshot
// @Security    BearerAuth
// @Router      /v1/user/device-data/{userDeviceID}/status [get]
func (d *DeviceDataController) GetUserDeviceStatus(c *fiber.Ctx) error {
	userDeviceID := c.Params("userDeviceID")

	userDevice, err := d.deviceAPI.GetUserDevice(c.Context(), userDeviceID)
	if err != nil {
		return err
	}

	deviceData, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(userDevice.Id),
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GT(time.Now().Add(-90*24*time.Hour)),
	).All(c.Context(), d.dbs().Reader)
	if err != nil {
		return err
	}

	if len(deviceData) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No status updates yet.")
	}

	ds := d.deviceStatusSvc.PrepareDeviceStatusInformation(c.Context(), deviceData, userDevice.DeviceDefinitionId,
		userDevice.DeviceStyleId, []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleCurrentLocation, privileges.VehicleAllTimeLocation})

	return c.JSON(ds)
}

// GetVehicleStatus godoc
// @Description Returns the latest status update for the vehicle with a given token id.
// @Tags        device-data
// @Param       tokenId path int true "token id"
// @Produce     json
// @Success     200 {object} response.DeviceSnapshot
// @Failure     404
// @Security    BearerAuth
// @Router      /v1/vehicle/{tokenId}/status [get]
func (d *DeviceDataController) GetVehicleStatus(c *fiber.Ctx) error {
	tis := c.Params("tokenID")
	claims := c.Locals("tokenClaims").(pr.CustomClaims)

	ti, ok := new(big.Int).SetString(tis, 10)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Couldn't parse token id %q.", tis))
	}

	// tid := pgtypes.NewNullDecimal(new(decimal.Big).SetBigMantScale(ti, 0))
	userDeviceNFT, err := d.deviceAPI.GetUserDeviceByTokenID(c.Context(), ti.Int64())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "NFT not found.")
		}
		d.log.Err(err).Str("token_id", tis).Msg("Database error retrieving NFT metadata or NFT not found")
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
	if errors.Is(err, sql.ErrNoRows) || len(deviceData) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no status updates yet")
	}
	if err != nil {
		return err
	}

	ds := d.deviceStatusSvc.PrepareDeviceStatusInformation(c.Context(), deviceData, userDeviceNFT.DeviceDefinitionId,
		userDeviceNFT.DeviceStyleId, claims.PrivilegeIDs)

	return c.JSON(ds)
}

// GetVehicleStatusRaw godoc
// @Description Returns the latest status update for the vehicle with a given token id.
// @Tags        device-data
// @Param       tokenId path int true "token id"
// @Produce     json
// @Success     200
// @Failure     404
// @Failure     500
// @Security    BearerAuth
// @Router      /v1/vehicle/{tokenId}/status-raw [get]
func (d *DeviceDataController) GetVehicleStatusRaw(c *fiber.Ctx) error {
	tis := c.Params("tokenID")
	privs := getPrivileges(c)

	ti, ok := new(big.Int).SetString(tis, 10)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Couldn't parse token id %q.", tis))
	}

	userDeviceNFT, err := d.deviceAPI.GetUserDeviceByTokenID(c.Context(), ti.Int64())
	if err != nil {
		d.log.Err(err).Msg("grpc error retrieving NFT metadata.")
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
	if errors.Is(err, sql.ErrNoRows) || len(deviceData) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no status updates yet")
	}
	if err != nil {
		return err
	}

	var jsonSignal []byte
	for _, item := range deviceData {
		jsonSignal = dataComplianceForSignals(item.Signals.JSON, privs)
		break
	}
	return c.Send(jsonSignal)
}

// dataComplianceForSignals removes any signals that per the privileges should not be returned
func dataComplianceForSignals(json []byte, privilegeIDs []privileges.Privilege) []byte {
	nonLocationDataProperties := []string{"charging", "fuelPercentRemaining", "batteryCapacity", "oil", "soc", "chargeLimit", "odometer", "range", "batteryVoltage", "ambientTemp", "tires"}
	currentLocationAndAllTimeLocation := []string{"latitude", "longitude"}

	removeProperties := func(properties []string) {
		for _, property := range properties {
			if result := gjson.GetBytes(json, property); result.Exists() {
				json, _ = sjson.DeleteBytes(json, property)
			}
		}
	}

	if !slices.Contains(privilegeIDs, privileges.VehicleNonLocationData) {
		removeProperties(nonLocationDataProperties)
	}

	if !slices.Contains(privilegeIDs, privileges.VehicleCurrentLocation) || slices.Contains(privilegeIDs, privileges.VehicleAllTimeLocation) {
		removeProperties(currentLocationAndAllTimeLocation)
	}

	return json
}

func multNull(x null.Float64, y float64) null.Float64 {
	if !x.Valid {
		return x
	}
	return null.Float64From(x.Float64 * y)
}

type odomValue struct {
	Value *float64 `json:"value"`
}

type dailyDistanceElasticResult struct {
	Aggregations struct {
		Days struct {
			Buckets []struct {
				KeyAsString string    `json:"key_as_string"`
				MinOdom     odomValue `json:"min_odom"`
				MaxOdom     odomValue `json:"max_odom"`
			} `json:"buckets"`
		} `json:"days"`
	} `json:"aggregations"`
}

type DailyDistanceDay struct {
	Date     string   `json:"date"`
	Distance *float64 `json:"distance"`
}

type DailyDistanceResp struct {
	Days []DailyDistanceDay `json:"days"`
}

// GetDailyDistance godoc
// @Description  [🔴__Warning - API Shutdown by June 30, 2024, Use `/v2/vehicles/:tokenId/analytics/daily-distance` instead__🔴] Get kilometers driven for a userDeviceID each day.
// @Tags         device-data [End Of Life Warning]
// @Produce      json
// @Success      200 {object} controllers.DailyDistanceResp
// @Failure      404 "no device found for user with provided parameters"
// @Param        userDeviceID  path   string  true   "user device id"
// @Param	     time_zone query string true "IANAS time zone id, e.g., America/Los_Angeles"
// @Security     BearerAuth
// @Router       /v1/user/device-data/{userDeviceID}/daily-distance [get]
func (d *DeviceDataController) GetDailyDistance(c *fiber.Ctx) error {
	userDeviceID := c.Params("userDeviceID")

	tz := c.Query("time_zone")

	query := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Match: map[string]types.MatchQuery{"subject": {Query: userDeviceID}}},
					{Range: map[string]types.RangeQuery{"data.timestamp": types.DateRangeQuery{Gte: some.String("now-13d/d"), TimeZone: &tz}}},
				},
			},
		},
		Size: some.Int(0),
		Aggregations: map[string]types.Aggregations{
			"days": {
				DateHistogram: &types.DateHistogramAggregation{
					Field:            some.String("data.timestamp"),
					CalendarInterval: &calendarinterval.Day,
					TimeZone:         &tz,
				},
				Aggregations: map[string]types.Aggregations{
					"min_odom": {
						Min: &types.MinAggregation{
							Field: some.String("data.odometer"),
						},
					},
					"max_odom": {
						Max: &types.MaxAggregation{
							Field: some.String("data.odometer"),
						},
					},
					// Code generation for buckets_path is broken as of 8.5.0
					// See https://github.com/elastic/go-elasticsearch/issues/570
				},
			},
		},
	}

	resp, err := d.esService.ESClient().Search().Index(d.Settings.DeviceDataIndexName).Request(query).Perform(c.Context())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c := resp.StatusCode; c != 200 {
		d.log.Error().Int("statusCode", c).Msg("Failed to get daily distance from Elastic.")
		// TODO(elffjs): Be more discerning here.
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}

	var ddr dailyDistanceElasticResult

	err = json.NewDecoder(resp.Body).Decode(&ddr)
	if err != nil {
		return err
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

// GetLastSeen godoc
// @Description  Specific for AutoPi - get when a device last sent data
// @Tags         autopi
// @Produce      json
// @Success      200
// @Failure      404 "no device found with eth addr or no data found"
// @Failure      400 "invalid eth addr"
// @Failure      500 "no device found or no data found, or other transient error"
// @Param        ethAddr  path   string  true   "device ethereum address"
// @Security     PreSharedKey
// @Router       /v1/autopi/last-seen/{ethAddr} [get]
func (d *DeviceDataController) GetLastSeen(c *fiber.Ctx) error {
	authed := false
	for h, vals := range c.GetReqHeaders() {
		if strings.EqualFold("Authorization", h) {
			for _, val := range vals {
				if val == d.Settings.AutoPiPreSharedKey {
					authed = true
					break
				}
			}
		}
	}
	if !authed {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	ethAddr := c.Params("ethAddr")
	if ethAddr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "empty eth addr")
	}
	if ethAddr[:2] != "0x" {
		return fiber.NewError(fiber.StatusBadRequest, "no valid eth addr found in route, must start with 0x")
	}
	addr := common.HexToAddress(ethAddr)

	ud, err := d.deviceAPI.GetUserDeviceByEthAddr(c.Context(), addr.Bytes())
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("aftermarket device not found with ethAddr: %s . %s", ethAddr, err.Error()))
	}
	udd, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(ud.Id)).One(c.Context(), d.dbs().Reader)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusFailedDependency, "device found, but no data reported: "+err.Error())
		}
		return err
	}
	return c.JSON(fiber.Map{
		"last_data_seen": udd.UpdatedAt.Format(time.RFC3339),
	})
}

func (d *DeviceDataController) GetDeviceDefinitionRawData(c *fiber.Ctx) error {
	//TODO: This is a temporary endpoint to get the raw data for a device definition
	userDeviceID := c.Params("userDeviceID", "")

	deviceData, err := models.UserDeviceData(models.UserDeviceDatumWhere.UserDeviceID.EQ(userDeviceID)).One(c.Context(), d.dbs().Reader)

	if err != nil {
		return err
	}

	d.log.Info().Interface("signals", deviceData.Signals).Msg("Returning raw data for device definition")

	return c.JSON(fiber.Map{
		"raw_data": deviceData.Signals,
	})
}
