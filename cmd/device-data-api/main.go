package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/DIMO-Network/device-data-api/docs"
	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/controllers"
	"github.com/DIMO-Network/device-data-api/internal/middleware/owner"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/middleware/privilegetoken"
	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

// @title                       DIMO Device Data API
// @version                     1.0
// @BasePath                    /v1
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
func main() {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("app", "device-data-api").
		Logger()

	settings, err := shared.LoadConfig[config.Settings]("settings.yaml")
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load settings")
	}
	level, err := zerolog.ParseLevel(settings.LogLevel)
	if err != nil {
		logger.Fatal().Err(err).Msgf("could not parse LOG_LEVEL: %s", settings.LogLevel)
	}
	zerolog.SetGlobalLevel(level)

	// start the actual stuff
	startPrometheus(logger)
	startWebAPI(logger, &settings)
}

func startWebAPI(logger zerolog.Logger, settings *config.Settings) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return ErrorHandler(c, err, logger)
		},
		DisableStartupMessage: true,
		ReadBufferSize:        16000,
	})

	app.Use(recover.New(recover.Config{
		Next:              nil,
		EnableStackTrace:  true,
		StackTraceHandler: nil,
	}))
	app.Use(cors.New())

	app.Get("/", healthCheck)
	app.Get("/v1/swagger/*", swagger.HandlerDefault)

	// secured paths
	keyRefreshInterval := time.Hour
	keyRefreshUnknownKID := true
	jwtAuth := jwtware.New(jwtware.Config{
		KeySetURL:            settings.JwtKeySetURL,
		KeyRefreshInterval:   &keyRefreshInterval,
		KeyRefreshUnknownKID: &keyRefreshUnknownKID,
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
	// establish grpc connections
	definitionsConn, err := grpc.Dial(settings.DeviceDefinitionsGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to dial device definitions grpc")
	}
	defer definitionsConn.Close()
	devicesConn, err := grpc.Dial(settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to dial devices grpc")
	}
	defer devicesConn.Close()

	usersConn, err := grpc.Dial(settings.UsersAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to dial users-api at %s", settings.UsersAPIGRPCAddr)
	}
	defer usersConn.Close()

	deviceAPIService := services.NewDeviceAPIService(devicesConn)
	definitionsAPIService := services.NewDeviceDefinitionsAPIService(definitionsConn)
	usersClient := pb.NewUserServiceClient(usersConn)

	deviceDataController := controllers.NewDeviceDataController(settings, &logger, deviceAPIService, esClient8, definitionsAPIService)

	logger.Info().Str("jwkUrl", settings.TokenExchangeJWTKeySetURL).Str("vehicleAddr", settings.VehicleNFTAddress).Msg("Privileges enabled.")
	privilegeAuth := jwtware.New(jwtware.Config{
		KeySetURL:            settings.TokenExchangeJWTKeySetURL,
		KeyRefreshInterval:   &keyRefreshInterval,
		KeyRefreshUnknownKID: &keyRefreshUnknownKID,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Err(err).Msg("Privilege token error.")
			return fiber.DefaultErrorHandler(c, err)
		},
	})

	vToken := app.Group("/v1/vehicle/:tokenID", privilegeAuth)

	tk := privilegetoken.New(privilegetoken.Config{
		Log: &logger,
	})
	vehicleAddr := common.HexToAddress(settings.VehicleNFTAddress)

	// Probably want constants for 1 and 4 here.
	vToken.Get("/history", tk.OneOf(vehicleAddr, []int64{controllers.NonLocationData, controllers.AllTimeLocation}), deviceDataController.GetHistoricalRawPermissioned)

	v1Auth := app.Group("/v1", jwtAuth)

	udMw := owner.New(usersClient, deviceAPIService, &logger)
	udOwner := v1Auth.Group("/user/device-data/:userDeviceID", udMw)
	udOwner.Get("/historical", deviceDataController.GetHistoricalRaw)
	udOwner.Get("/distance-driven", deviceDataController.GetDistanceDriven)
	udOwner.Get("/daily-distance", deviceDataController.GetDailyDistance)

	if settings.Environment != "prod" {
		dataDownloadController, err := controllers.NewDataDownloadController(settings, &logger, esClient8, deviceAPIService)
		if err != nil {
			panic(err)
		}

		udOwner.Get("/export/json/email", dataDownloadController.DataDownloadHandler)
		go func() {
			err = dataDownloadController.DataDownloadConsumer(context.Background())
			if err != nil {
				logger.Info().Err(err).Msg("data download consumer error")
			}
		}()
	}

	logger.Info().Msg("Server started on port " + settings.Port)
	// Start Server from a different go routine
	go func() {
		if err := app.Listen(":" + settings.Port); err != nil {
			logger.Fatal().Err(err)
		}
	}()
	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent with length of 1
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-c                                             // This blocks the main thread until an interrupt is received
	logger.Info().Msg("Gracefully shutting down and running cleanup tasks...")
	_ = app.Shutdown()
	// shutdown anything else
}

func startPrometheus(logger zerolog.Logger) {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":8888", nil)
		if err != nil {
			logger.Fatal().Err(err).Msg("could not start consumer")
		}
	}()
	logger.Info().Msg("prometheus metrics at :8888/metrics")
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

	logger.Err(err).Int("code", code).Str("path", strings.TrimPrefix(c.Path(), "/")).Msg("Failed request.")

	return c.Status(code).JSON(CodeResp{Code: code, Message: message})
}

type CodeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(CodeResp{Code: 200, Message: "Server is up."})
}
