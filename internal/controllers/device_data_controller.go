package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared/db"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/DIMO-Network/shared"

	"github.com/DIMO-Network/devices-api/pkg/grpc"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	pr "github.com/DIMO-Network/shared/middleware/privilegetoken"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/exp/slices"
)

type DeviceDataController struct {
	Settings       *config.Settings
	log            *zerolog.Logger
	deviceAPI      services.DeviceAPIService
	es8Client      *es8.TypedClient
	definitionsAPI services.DeviceDefinitionsAPIService
	dbs            func() *db.ReaderWriter
}

const (
	NonLocationData int64 = 1
	Commands        int64 = 2
	CurrentLocation int64 = 3
	AllTimeLocation int64 = 4
)

// NewDeviceDataController constructor
func NewDeviceDataController(settings *config.Settings, logger *zerolog.Logger, deviceAPIService services.DeviceAPIService, es8Client *es8.TypedClient, definitionsAPIService services.DeviceDefinitionsAPIService, dbs func() *db.ReaderWriter) DeviceDataController {
	return DeviceDataController{
		Settings:       settings,
		log:            logger,
		deviceAPI:      deviceAPIService,
		es8Client:      es8Client,
		definitionsAPI: definitionsAPIService,
		dbs:            dbs,
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
// @Router       /user/device-data/{userDeviceID}/historical [get]
func (d *DeviceDataController) GetHistoricalRaw(c *fiber.Ctx) error {
	const dateLayout = "2006-01-02" // date layout support by elastic
	userDeviceID := c.Params("userDeviceID")
	startDate := c.Query("startDate")
	if startDate == "" {
		startDate = time.Now().Add(-1 * (time.Hour * 24 * 14)).Format(dateLayout)
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

	return d.getHistory(c, userDevice, startDate, endDate, types.SourceFilter{})
}

// addRangeIfNotExists will add range based on mpg and fuelTankCapacity to the json body, only if there are no existing range entries (eg. smartcar added)
func addRangeIfNotExists(ctx context.Context, deviceDefSvc services.DeviceDefinitionsAPIService, body []byte, deviceDefinitionID string, deviceStyleID *string) ([]byte, error) {
	if len(body) == 0 {
		return body, nil
	}
	// check if range is already present in any document
	if gjson.GetBytes(body, "hits.hits.#(_source.data.range>0)0._source.data.range").Exists() {
		return body, nil
	}

	definition, err := deviceDefSvc.GetDeviceDefinitionByID(ctx, deviceDefinitionID)
	if err != nil {
		return body, errors.Wrapf(err, "could not get device definition by id: %s", deviceDefinitionID)
	}
	// extract the range values from definition, already done in devices-api, copy that code or move to shared
	rangeData := GetActualDeviceDefinitionMetadataValues(definition, deviceStyleID)

	resultData := gjson.GetBytes(body, "hits.hits.#._source.data")
	for i, r := range resultData.Array() {
		// note range is reported in km
		fuelResult := r.Get("fuelPercentRemaining")
		if fuelResult.Exists() {
			rangeKm := CalculateRange(rangeData, fuelResult.Num)
			if rangeKm != nil {
				body, err = sjson.SetBytes(body, fmt.Sprintf("hits.hits.%d._source.data.range", i), rangeKm)
				if err != nil {
					return body, err
				}
			}
		}
	}

	return body, nil
}

// GetHistoricalRawPermissioned godoc
// @Description  Get all historical data for a tokenID, within start and end range
// @Tags         device-data
// @Produce      json
// @Success      200
// @Param        tokenID  path   int64  true   "token id"
// @Param        startDate     query  string  false  "startDate eg 2022-01-02. if empty two weeks back"
// @Param        endDate       query  string  false  "endDate eg 2022-03-01. if empty today"
// @Security     BearerAuth
// @Router       /vehicle/{tokenID}/history [get]
func (d *DeviceDataController) GetHistoricalRawPermissioned(c *fiber.Ctx) error {
	const dateLayout = "2006-01-02" // date layout support by elastic
	tokenID := c.Params("tokenID")
	startDate := c.Query("startDate")
	if startDate == "" {
		startDate = time.Now().Add(-1 * (time.Hour * 24 * 14)).Format(dateLayout)
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

	i, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userDevice, err := d.deviceAPI.GetUserDeviceByTokenID(c.Context(), i)
	if err != nil {
		return err
	}

	claims := c.Locals("tokenClaims").(pr.CustomClaims)
	privileges := claims.PrivilegeIDs

	var filter types.SourceFilter

	if slices.Contains(privileges, AllTimeLocation) {
		filter.Includes = append(filter.Includes, "data.latitude", "data.longitude", "location", "data.cell", "cell")
	} else {
		filter.Excludes = append(filter.Excludes, "data.latitude", "data.longitude", "location", "data.cell", "cell")
	}

	if slices.Contains(privileges, NonLocationData) {
		// Overrides the more limited Includes entries from above if the token also
		// has AllTimeLocation.
		filter.Includes = append(filter.Includes, "*")
	}

	return d.getHistory(c, userDevice, startDate, endDate, filter)
}

func (d *DeviceDataController) getHistory(c *fiber.Ctx, userDevice *grpc.UserDevice, startDate, endDate string, filter types.SourceFilter) error {
	msm := types.MinimumShouldMatch(1)

	var source types.SourceConfig = filter

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

	res, err := d.es8Client.Search().Index(d.Settings.DeviceDataIndexName).Request(&req).Do(c.Context())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	localLog := d.log.With().Str("userDeviceId", userDevice.Id).
		Str("deviceDefinitionId", userDevice.DeviceDefinitionId).Interface("response", res).Logger()
	if res.StatusCode >= fiber.StatusBadRequest {
		localLog.Error().Str("userDeviceId", userDevice.Id).Interface("response", res).Msgf("Got status code %d from Elastic.", res.StatusCode)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		localLog.Err(err).Str("userDeviceId", userDevice.Id).Msg("Failed to read Elastic response body.")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}

	body, err = addRangeIfNotExists(c.Context(), d.definitionsAPI, body, userDevice.DeviceDefinitionId, userDevice.DeviceStyleId)
	if err != nil {
		localLog.Warn().Err(err).Msg("could not add range calculation to document")
	}
	body = removeOdometerIfInvalid(body)

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

// GetDistanceDriven godoc
// @Description  Get kilometers driven for a userDeviceID since connected (ie. since we have data available)
// @Description  if it returns 0 for distanceDriven it means we have no odometer data.
// @Tags         device-data
// @Produce      json
// @Success      200
// @Failure      404 "no device found for user with provided parameters"
// @Param        userDeviceID  path   string  true   "user device id"
// @Security     BearerAuth
// @Router       /user/device-data/{userDeviceID}/distance-driven [get]
func (d *DeviceDataController) GetDistanceDriven(c *fiber.Ctx) error {
	userDeviceID := c.Params("userDeviceID")

	odoStart, err := d.queryOdometer(c.Context(), sortorder.Asc, userDeviceID)
	if err != nil {
		return errors.Wrap(err, "error querying odometer")
	}
	odoEnd, err := d.queryOdometer(c.Context(), sortorder.Desc, userDeviceID)
	if err != nil {
		return errors.Wrap(err, "error querying odometer")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"distanceDriven": odoEnd - odoStart,
		"units":          "kilometers",
	})
}

// GetUserDeviceStatus godoc
// @Description Returns the latest status update for the device. May return 404 if the
// @Description user does not have a device with the ID, or if no status updates have come. Note this endpoint also exists under nft_controllers
// @Tags        user-devices
// @Produce     json
// @Param       user_device_id path     string true "user device ID"
// @Success     200            {object} controllers.DeviceSnapshot
// @Security    BearerAuth
// @Router      /user/device-data/{userDeviceID}/status [get]
func (d *DeviceDataController) GetUserDeviceStatus(c *fiber.Ctx) error {
	userDeviceID := c.Params("userDeviceID")

	userDevice, err := d.deviceAPI.GetUserDevice(c.Context(), userDeviceID)
	if err != nil {
		return err
	}

	deviceData, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(userDevice.Id),
		models.UserDeviceDatumWhere.Signals.IsNotNull(),
		models.UserDeviceDatumWhere.UpdatedAt.GT(time.Now().Add(-14*24*time.Hour)),
	).All(c.Context(), d.dbs().Reader)
	if err != nil {
		return err
	}

	if len(deviceData) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No status updates yet.")
	}

	ds := PrepareDeviceStatusInformation(c.Context(), d.definitionsAPI, deviceData, userDevice.DeviceDefinitionId,
		userDevice.DeviceStyleId, []int64{NonLocationData, CurrentLocation, AllTimeLocation})

	return c.JSON(ds)
}

// GetVehicleStatus godoc
// @Description Returns the latest status update for the vehicle with a given token id.
// @Tags        permission
// @Param       tokenId path int true "token id"
// @Produce     json
// @Success     200 {object} controllers.DeviceSnapshot
// @Failure     404
// @Router      /vehicle/{tokenId}/status [get]
func (d *DeviceDataController) GetVehicleStatus(c *fiber.Ctx) error {
	tis := c.Params("tokenID")
	claims := c.Locals("tokenClaims").(pr.CustomClaims)

	privileges := claims.PrivilegeIDs

	ti, ok := new(big.Int).SetString(tis, 10)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Couldn't parse token id %q.", tis))
	}

	//tid := pgtypes.NewNullDecimal(new(decimal.Big).SetBigMantScale(ti, 0))
	userDeviceNFT, err := d.deviceAPI.GetUserDeviceByTokenID(c.Context(), ti.Int64())
	if err != nil {
		d.log.Err(err).Msg("grpc error retrieving NFT metadata.")
		return err
	}

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
		models.UserDeviceDatumWhere.UpdatedAt.GT(time.Now().Add(-14*24*time.Hour)),
	).All(c.Context(), d.dbs().Reader)
	if errors.Is(err, sql.ErrNoRows) || len(deviceData) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no status updates yet")
	}
	if err != nil {
		return err
	}

	ds := PrepareDeviceStatusInformation(c.Context(), d.definitionsAPI, deviceData, userDeviceNFT.DeviceDefinitionId,
		userDeviceNFT.DeviceStyleId, privileges)

	return c.JSON(ds)
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
// @Description  Get kilometers driven for a userDeviceID each day.
// @Tags         device-data
// @Produce      json
// @Success      200 {object} controllers.DailyDistanceResp
// @Failure      404 "no device found for user with provided parameters"
// @Param        userDeviceID  path   string  true   "user device id"
// @Param	     time_zone query string true "IANAS time zone id, e.g., America/Los_Angeles"
// @Security     BearerAuth
// @Router       /user/device-data/{userDeviceID}/daily-distance [get]
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

	resp, err := d.es8Client.Search().Index(d.Settings.DeviceDataIndexName).Request(query).Do(c.Context())
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

// queryOdometer gets the lowest or highest odometer reading depending on order - asc = lowest, desc = highest
func (d *DeviceDataController) queryOdometer(ctx context.Context, order sortorder.SortOrder, userDeviceID string) (float64, error) {
	req := search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{Term: map[string]types.TermQuery{"subject": {Value: userDeviceID}}},
					{Exists: &types.ExistsQuery{Field: "data.odometer"}},
				},
			},
		},
		Size: some.Int(1),
		Sort: []types.SortCombinations{types.SortOptions{SortOptions: map[string]types.FieldSort{"data.odometer": {Order: &order}}}},
	}

	res, err := d.es8Client.Search().Index(d.Settings.DeviceDataIndexName).Request(&req).Do(ctx)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, fmt.Errorf("status code %d from Elastic", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	if gjson.GetBytes(body, "hits.hits.#").Int() == 0 {
		// Existing behavior. Not great.
		return 0, nil
	}

	body = removeOdometerIfInvalid(body)

	return gjson.GetBytes(body, "hits.hits.0._source.data.odometer").Float(), nil
}
