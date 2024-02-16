package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/controllers"
	"github.com/DIMO-Network/device-data-api/internal/middleware/metrics"
	"github.com/DIMO-Network/device-data-api/internal/middleware/owner"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/device-data-api/internal/services/elastic"
	"github.com/DIMO-Network/shared/db"
	"github.com/DIMO-Network/shared/middleware/privilegetoken"
	"github.com/DIMO-Network/shared/privileges"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/ethereum/go-ethereum/common"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/skip"
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

	cacheHandler := skip.New(cache.New(cache.Config{
		Expiration:   2 * time.Minute,
		CacheControl: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.OriginalURL() + c.Get(fiber.HeaderAuthorization)
		},
	}), func(c *fiber.Ctx) bool {
		// skip cache if refresh is true
		return c.Query("refresh") == "true"
	})

	app.Get("/", healthCheck)
	app.Get("/v1/swagger/*", swagger.HandlerDefault)

	// secured paths
	jwtAuth := jwtware.New(jwtware.Config{
		JWKSetURLs: []string{settings.JwtKeySetURL},
		ErrorHandler: func(_ *fiber.Ctx, _ error) error {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT.")
		},
	})

	esService, err := elastic.New(settings, &logger, nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error constructing Elasticsearch service.")
	}
	deviceStatusSvc := services.NewDeviceStatusService(definitionsAPIService)
	deviceDataController := controllers.NewDeviceDataController(settings, &logger, deviceAPIService, esService, definitionsAPIService, deviceStatusSvc, dbs)

	logger.Info().Str("jwkUrl", settings.TokenExchangeJWTKeySetURL).Str("vehicleAddr", settings.VehicleNFTAddress).Msg("Privileges enabled.")
	privilegeAuth := jwtware.New(jwtware.Config{
		JWKSetURLs: []string{settings.TokenExchangeJWTKeySetURL},
		ErrorHandler: func(_ *fiber.Ctx, err error) error {
			logger.Err(err).Msg("Privilege token error.")
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid privilege token.")
		},
	})
	// autopi specific endpoint,
	app.Get("/v1/autopi/last-seen/:ethAddr", deviceDataController.GetLastSeen)

	vTokenV1 := app.Group("/v1/vehicle/:tokenID", privilegeAuth)
	vTokenV2 := app.Group("/v2/vehicle/:tokenID", privilegeAuth)

	tk := privilegetoken.New(privilegetoken.Config{
		Log: &logger,
	})
	vehicleAddr := common.HexToAddress(settings.VehicleNFTAddress)

	// token based routes
	vTokenV1.Get("/history", tk.OneOf(vehicleAddr, []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation}), cacheHandler, deviceDataController.GetHistoricalRawPermissioned)
	vTokenV1.Get("/status", tk.OneOf(vehicleAddr, []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleCurrentLocation, privileges.VehicleAllTimeLocation}), cacheHandler, deviceDataController.GetVehicleStatus)
	vTokenV1.Get("/status-raw", tk.OneOf(vehicleAddr, []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleCurrentLocation, privileges.VehicleAllTimeLocation}), cacheHandler, deviceDataController.GetVehicleStatusRaw)

	vTokenV2.Get("/status", tk.OneOf(vehicleAddr, []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleCurrentLocation, privileges.VehicleAllTimeLocation}), cacheHandler, deviceDataController.GetVehicleStatusV2)
	vTokenV2.Get("/history", tk.OneOf(vehicleAddr, []privileges.Privilege{privileges.VehicleNonLocationData, privileges.VehicleAllTimeLocation}), cacheHandler, deviceDataController.GetHistoricalPermissionedV2)

	udMw := owner.New(usersClient, deviceAPIService, &logger)
	v1Auth := app.Group("/v1", jwtAuth)

	udOwner := v1Auth.Group("/user/device-data/:userDeviceID", udMw)
	udOwner.Get("/status", cacheHandler, deviceDataController.GetUserDeviceStatus)
	udOwner.Get("/historical", cacheHandler, deviceDataController.GetHistoricalRaw)
	udOwner.Get("/distance-driven", cacheHandler, deviceDataController.GetDistanceDriven)
	udOwner.Get("/daily-distance", cacheHandler, deviceDataController.GetDailyDistance)

	dataDownloadController, err := controllers.NewDataDownloadController(settings, &logger, esService.ESClient(), deviceAPIService)
	if err != nil {
		panic(err)
	}

	udOwner.Post("/export/json/email", dataDownloadController.DataDownloadHandler)
	go func() {
		err = dataDownloadController.DataDownloadConsumer(context.Background())
		if err != nil {
			logger.Err(err).Msg("data download consumer error")
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

// ErrorHandler custom handler to log recovered errors using our logger and return json instead of string
func ErrorHandler(c *fiber.Ctx, err error, logger zerolog.Logger) error {
	code := fiber.StatusInternalServerError // Default 500 statuscode
	message := "Internal error."

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}
	if strings.Contains(err.Error(), "code = NotFound") {
		code = fiber.StatusNotFound
	}
	const descPrefix = "desc = "
	if strings.Contains(err.Error(), descPrefix) {
		// pull out desc from message - typically grpc errors
		start := strings.Index(err.Error(), descPrefix)
		start += len(descPrefix)
		message = err.Error()[start:]
	}

	// don't log not found errors
	if code != fiber.StatusNotFound {
		logger.Err(err).Int("httpStatusCode", code).
			Str("httpPath", strings.TrimPrefix(c.Path(), "/")).
			Str("httpMethod", c.Method()).
			Msg("caught an error from http request")
	}

	return c.Status(code).JSON(CodeResp{Code: code, Message: message})
}

type CodeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
