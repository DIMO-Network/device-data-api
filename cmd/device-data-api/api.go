package main

import (
	"context"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/controllers"
	"github.com/DIMO-Network/device-data-api/internal/middleware/metrics"
	"github.com/DIMO-Network/device-data-api/internal/middleware/owner"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/shared/db"
	"github.com/DIMO-Network/shared/middleware/privilegetoken"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ethereum/go-ethereum/common"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog"
)

func startWebAPI(logger zerolog.Logger, settings *config.Settings, dbs func() *db.ReaderWriter,
	definitionsAPIService services.DeviceDefinitionsAPIService,
	deviceAPIService services.DeviceAPIService,
	usersClient pb.UserServiceClient) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return ErrorHandler(c, err, logger)
		},
		DisableStartupMessage: true,
		ReadBufferSize:        16000,
	})

	app.Use(metrics.HTTPMetricsMiddleware)

	app.Use(fiberrecover.New(fiberrecover.Config{
		Next:              nil,
		EnableStackTrace:  true,
		StackTraceHandler: nil,
	}))
	app.Use(cors.New())
	//cache
	cacheHandler := cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   2 * time.Minute,
		CacheControl: true,
	})

	app.Get("/", healthCheck)
	app.Get("/v1/swagger/*", swagger.HandlerDefault)

	// secured paths
	jwtAuth := jwtware.New(jwtware.Config{
		JWKSetURLs: []string{settings.JwtKeySetURL},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT.")
		},
	})

	esClient8, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses:  []string{settings.ElasticSearchAnalyticsHost},
		Username:   settings.ElasticSearchAnalyticsUsername,
		Password:   settings.ElasticSearchAnalyticsPassword,
		MaxRetries: 5,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Error constructing Elasticsearch client.")
	}

	deviceStatusSvc := services.NewDeviceStatusService(definitionsAPIService)

	deviceDataController := controllers.NewDeviceDataController(settings, &logger, deviceAPIService, esClient8, definitionsAPIService, deviceStatusSvc, dbs)

	logger.Info().Str("jwkUrl", settings.TokenExchangeJWTKeySetURL).Str("vehicleAddr", settings.VehicleNFTAddress).Msg("Privileges enabled.")
	privilegeAuth := jwtware.New(jwtware.Config{
		JWKSetURLs: []string{settings.TokenExchangeJWTKeySetURL},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Err(err).Msg("Privilege token error.")
			return fiber.DefaultErrorHandler(c, err)
		},
	})
	// autopi specific endpoint,
	app.Get("/v1/autopi/last-seen/:ethAddr", cacheHandler, deviceDataController.GetLastSeen)

	vToken := app.Group("/v1/vehicle/:tokenID", privilegeAuth)

	tk := privilegetoken.New(privilegetoken.Config{
		Log: &logger,
	})
	vehicleAddr := common.HexToAddress(settings.VehicleNFTAddress)

	// token based routes
	vToken.Get("/history", tk.OneOf(vehicleAddr, []int64{controllers.NonLocationData, controllers.AllTimeLocation}), cacheHandler, deviceDataController.GetHistoricalRawPermissioned)
	vToken.Get("/status", tk.OneOf(vehicleAddr, []int64{controllers.NonLocationData, controllers.CurrentLocation, controllers.AllTimeLocation}), cacheHandler, deviceDataController.GetVehicleStatus)

	v1Auth := app.Group("/v1", jwtAuth)

	udMw := owner.New(usersClient, deviceAPIService, &logger)
	udOwner := v1Auth.Group("/user/device-data/:userDeviceID", udMw)
	udOwner.Get("/status", cacheHandler, deviceDataController.GetUserDeviceStatus)
	udOwner.Get("/historical", cacheHandler, deviceDataController.GetHistoricalRaw)
	udOwner.Get("/distance-driven", cacheHandler, deviceDataController.GetDistanceDriven)
	udOwner.Get("/daily-distance", cacheHandler, deviceDataController.GetDailyDistance)

	dataDownloadController, err := controllers.NewDataDownloadController(settings, &logger, esClient8, deviceAPIService)
	if err != nil {
		panic(err)
	}

	udOwner.Post("/export/json/email", dataDownloadController.DataDownloadHandler)
	go func() {
		err = dataDownloadController.DataDownloadConsumer(context.Background())
		if err != nil {
			logger.Info().Err(err).Msg("data download consumer error")
		}
	}()

	logger.Info().Msg("Server started on port " + settings.Port)
	// Start Server from a different go routine
	go func() {
		if err := app.Listen(":" + settings.Port); err != nil {
			logger.Fatal().Err(err).Send()
		}
	}()
	return app
}
