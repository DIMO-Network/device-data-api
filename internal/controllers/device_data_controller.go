package controllers

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

type DeviceDataController struct {
	Settings  *config.Settings
	log       *zerolog.Logger
	es        *elasticsearch.Client
	deviceAPI services.DeviceAPIService
}

// NewDeviceDataController constructor
func NewDeviceDataController(settings *config.Settings, logger *zerolog.Logger) DeviceDataController {
	es, err := connect(settings)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not connect to elastic search")
	}
	return DeviceDataController{
		Settings:  settings,
		log:       logger,
		es:        es,
		deviceAPI: services.NewDeviceAPIService(settings.DevicesAPIGRPCAddr),
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
	userID := getUserID(c)
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

	// todo: cache user devices in memory
	exists, err := d.deviceAPI.UserDeviceBelongsToUserID(c.Context(), userID, userDeviceID)
	if err != nil {
		return err
	}
	if !exists {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	res, err := esquery.Search().
		Query(
			esquery.CustomQuery(
				map[string]any{
					"function_score": map[string]any{
						"query": esquery.Bool().
							Filter(
								esquery.Term("subject", userDeviceID),
								esquery.Range("data.timestamp").Gte(startDate).Lte(endDate),
							).
							Should(
								esquery.Exists("data.odometer"),
								esquery.Exists("data.latitude"),
							).
							MinimumShouldMatch(1).Map(),
						"random_score": map[string]any{},
					},
				},
			),
		).
		Size(1000).
		Run(d.es, d.es.Search.WithContext(c.Context()), d.es.Search.WithIndex(d.Settings.DeviceDataIndexName))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= fiber.StatusBadRequest {
		d.log.Error().Str("userDeviceId", userDeviceID).Interface("response", res).Msgf("Got status code %d from Elastic.", res.StatusCode)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		d.log.Err(err).Str("userDeviceId", userDeviceID).Msg("Failed to read Elastic response body.")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error.")
	}
	c.Set("Content-Type", fiber.MIMEApplicationJSON)
	return c.Status(fiber.StatusOK).Send(body)
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
	userID := getUserID(c)
	userDeviceID := c.Params("userDeviceID")

	exists, err := d.deviceAPI.UserDeviceBelongsToUserID(c.Context(), userID, userDeviceID)
	if err != nil {
		return err
	}
	if !exists {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("No device %s found for user %s", userDeviceID, userID))
	}

	odoStart, err := d.queryOdometer(c.Context(), esquery.OrderAsc, userDeviceID)
	if err != nil {
		return errors.Wrap(err, "error querying odometer")
	}
	odoEnd, err := d.queryOdometer(c.Context(), esquery.OrderDesc, userDeviceID)
	if err != nil {
		return errors.Wrap(err, "error querying odometer")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"distanceDriven": odoEnd - odoStart,
		"units":          "kilometers",
	})
}

// queryOdometer gets the first or last odometer reading depending on order - asc = first, desc = last
func (d *DeviceDataController) queryOdometer(ctx context.Context, order esquery.Order, userDeviceID string) (float64, error) {
	res, err := esquery.Search().SourceIncludes("data.odometer").
		Query(esquery.Bool().Must(
			esquery.Term("subject", userDeviceID),
			esquery.Exists("data.odometer"),
		)).
		Size(1).
		Sort("data.timestamp", order).
		Run(d.es, d.es.Search.WithContext(ctx), d.es.Search.WithIndex(d.Settings.DeviceDataIndexName))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close() // nolint
	body, _ := io.ReadAll(res.Body)
	d.log.Info().RawJSON("respBody", body).Str("userDeviceID", userDeviceID).Str("order", string(order)).Msg("Queried for odometer")
	result := gjson.GetBytes(body, "hits.hits.0._source.data.odometer")
	if result.Exists() {
		return result.Float(), nil
	}
	return 0, nil
}

// connect helper to connect to ES. Move this to seperate file or under services etc.
func connect(settings *config.Settings) (*elasticsearch.Client, error) {
	// maybe refactor some of this into elasticsearchservice

	if settings.ElasticSearchAnalyticsUsername == "" {
		// we're connecting to local instance at localhost:9200
		return elasticsearch.NewDefaultClient()
	}

	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses:            []string{settings.ElasticSearchAnalyticsHost},
		Username:             settings.ElasticSearchAnalyticsUsername,
		Password:             settings.ElasticSearchAnalyticsPassword,
		EnableRetryOnTimeout: true,
		MaxRetries:           5,
	})
}
