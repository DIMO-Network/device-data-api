package controllers

import (
	"io"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/aquasecurity/esquery"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
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

	// todo: cache user devices in memeory
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
				map[string]interface{}{
					"function_score": map[string]interface{}{
						"query": esquery.Bool().
							Filter(
								esquery.Term("subject", userDeviceID),
								esquery.Range("data.timestamp").Gte(startDate).Lte(endDate),
							).
							Should(
								esquery.Exists("data.odometer"),
								esquery.Exists("data.latitude"),
							).
							MinimumShouldMatch(1),
						"random_score": map[string]interface{}{},
					},
				},
			),
		).
		Size(500).
		Sort("data.timestamp", "desc").
		Run(d.es, d.es.Search.WithContext(c.Context()), d.es.Search.WithIndex(d.Settings.DeviceDataIndexName))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	c.Set("Content-type", "application/json")
	return c.Status(fiber.StatusOK).Send(body)
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
