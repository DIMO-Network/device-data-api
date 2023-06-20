package main

import (
	"context"
	"errors"
	"github.com/DIMO-Network/device-data-api/internal/middleware/metrics"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/burdiyan/kafkautil"
	"github.com/google/subcommands"
	"github.com/lovoo/goka"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cache"

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
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	ctx := context.Background()
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

	pdb := db.NewDbConnectionFromSettings(ctx, &settings.DB, true)
	// check db ready, this is not ideal btw, the db connection handler would be nicer if it did this.
	totalTime := 0
	for !pdb.IsReady() {
		if totalTime > 30 {
			logger.Fatal().Msg("could not connect to postgres after 30 seconds")
		}
		time.Sleep(time.Second)
		totalTime++
	}

	deps := newDependencyContainer(&settings, logger, pdb.DBS)

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	// start the actual stuff
	if len(os.Args) == 1 {
		startPrometheus(logger)
		eventService := services.NewEventService(&logger, &settings, deps.getKafkaProducer())
		startDeviceStatusConsumer(logger, &settings, pdb, eventService, deps.getDeviceDefinitionService(), deps.getDeviceService())
		startWebAPI(logger, &settings, deps.getDeviceDefinitionService(), deps.getDeviceService())
	}
}

func startWebAPI(logger zerolog.Logger, settings *config.Settings, definitionsAPIService services.DeviceDefinitionsAPIService, deviceAPIService services.DeviceAPIService) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return ErrorHandler(c, err, logger)
		},
		DisableStartupMessage: true,
		ReadBufferSize:        16000,
	})

	app.Use(metrics.HTTPMetricsMiddleware)

	app.Use(recover.New(recover.Config{
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

	usersConn, err := grpc.Dial(settings.UsersAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal().Err(err).Msgf("Failed to dial users-api at %s", settings.UsersAPIGRPCAddr)
	}
	defer usersConn.Close()
	usersClient := pb.NewUserServiceClient(usersConn)

	deviceDataController := controllers.NewDeviceDataController(settings, &logger, deviceAPIService, esClient8, definitionsAPIService)

	logger.Info().Str("jwkUrl", settings.TokenExchangeJWTKeySetURL).Str("vehicleAddr", settings.VehicleNFTAddress).Msg("Privileges enabled.")
	privilegeAuth := jwtware.New(jwtware.Config{
		JWKSetURLs: []string{settings.TokenExchangeJWTKeySetURL},
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

func startDeviceStatusConsumer(logger zerolog.Logger, settings *config.Settings, pdb db.Store, eventService services.EventService,
	ddSvc services.DeviceDefinitionsAPIService, deviceSvc services.DeviceAPIService) {

	autoPISvc := services.NewAutoPiAPIService(settings, pdb.DBS)
	ingestSvc := services.NewDeviceStatusIngestService(pdb.DBS, &logger, eventService, ddSvc, autoPISvc, deviceSvc)

	sc := goka.DefaultConfig()
	sc.Version = sarama.V2_8_1_0
	goka.ReplaceGlobalConfig(sc)

	group := goka.DefineGroup("devices-data-consumer",
		goka.Input(goka.Stream(settings.DeviceStatusTopic), new(shared.JSONCodec[services.DeviceStatusEvent]), ingestSvc.ProcessDeviceStatusMessages),
	)

	processor, err := goka.NewProcessor(strings.Split(settings.KafkaBrokers, ","),
		group,
		goka.WithHasher(kafkautil.MurmurHasher),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not start device status processor")
	}

	go func() {
		err = processor.Run(context.Background())
		if err != nil {
			logger.Fatal().Err(err).Msg("could not run device status processor")
		}
	}()

	logger.Info().Msg("Device status update consumer started")
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
