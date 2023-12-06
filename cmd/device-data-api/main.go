package main

import (
	"context"
	"errors"
	"flag"
	"os/signal"
	"syscall"

	"github.com/DIMO-Network/device-data-api/internal/services/fingerprint"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	pb "github.com/DIMO-Network/users-api/pkg/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/DIMO-Network/device-data-api/internal/rpc"

	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/middleware/metrics"
	"github.com/DIMO-Network/shared/db"
	"github.com/Shopify/sarama"
	"github.com/burdiyan/kafkautil"
	"github.com/google/subcommands"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/lovoo/goka"

	_ "github.com/DIMO-Network/device-data-api/docs"
	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	dddatagrpc "github.com/DIMO-Network/device-data-api/pkg/grpc"
	"github.com/DIMO-Network/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// @title                       DIMO Device Data API
// @version                     1.0
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

	deviceDefsSvc, deviceDefsConn := deps.getDeviceDefinitionService()
	defer deviceDefsConn.Close()
	devicesSvc, devicesConn := deps.getDeviceService()
	defer devicesConn.Close()

	// start the actual stuff
	if len(os.Args) == 1 {
		startPrometheus(logger)

		go startGRPCServer(&settings, pdb.DBS, &logger, deviceDefsSvc)

		if settings.IsKafkaEnabled(&logger) {
			eventService := services.NewEventService(&logger, &settings, deps.getKafkaProducer())
			startDeviceStatusConsumer(logger, &settings, pdb, eventService, deviceDefsSvc, devicesSvc)
			startDeviceFingerprint(logger, &settings, pdb, devicesSvc)
		}
		if settings.IsWebAPIEnabled(&logger) {
			usersConn, err := grpc.Dial(settings.UsersAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				logger.Fatal().Err(err).Msgf("Failed to dial users-api at %s", settings.UsersAPIGRPCAddr)
			}
			defer usersConn.Close()
			usersClient := pb.NewUserServiceClient(usersConn)
			app := startWebAPI(logger, &settings, pdb.DBS, deviceDefsSvc, devicesSvc, usersClient)
			// nolint
			defer app.Shutdown()
		}

		c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent with length of 1
		signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
		<-c                                             // This blocks the main thread until an interrupt is received
		logger.Info().Msg("Gracefully shutting down and running cleanup tasks...")
		// shutdown anything else
	} else {
		subcommands.Register(&migrateDBCmd{logger: logger, settings: settings}, "database")
		subcommands.Register(&vehicleSignalsBatchCmd{db: pdb.DBS, logger: logger, deviceDefSvc: deviceDefsSvc, deviceSvc: devicesSvc}, "events")

		flag.Parse()
		os.Exit(int(subcommands.Execute(ctx)))
	}
}

func startGRPCServer(settings *config.Settings, dbs func() *db.ReaderWriter, logger *zerolog.Logger, definitionsAPIService services.DeviceDefinitionsAPIService) {
	lis, err := net.Listen("tcp", ":"+settings.GRPCPort)
	if err != nil {
		logger.Fatal().Err(err).Msgf("Couldn't listen on gRPC port %s", settings.GRPCPort)
	}

	logger.Info().Msgf("Starting gRPC server on port %s", settings.GRPCPort)
	gp := metrics.GRPCPanicker{Logger: logger}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			metrics.GRPCMetricsAndLogMiddleware(logger),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(gp.GRPCPanicRecoveryHandler)),
		)),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)

	deviceStatusSvc := services.NewDeviceStatusService(definitionsAPIService)
	dddatagrpc.RegisterUserDeviceDataServiceServer(server, rpc.NewUserDeviceData(dbs, logger, definitionsAPIService, deviceStatusSvc))

	if err := server.Serve(lis); err != nil {
		logger.Fatal().Err(err).Msg("gRPC server terminated unexpectedly")
	}
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

func startDeviceFingerprint(logger zerolog.Logger, settings *config.Settings, pdb db.Store, deviceAPIService services.DeviceAPIService) {
	ctx := context.Background()

	if err := fingerprint.RunConsumer(ctx, settings, &logger, pdb, deviceAPIService); err != nil {
		logger.Fatal().Err(err).Msg("Failed to create vin credentialer listener")
	}
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

	// don't log not found errors
	if code != fiber.StatusNotFound {
		logger.Err(err).Int("code", code).Str("path", strings.TrimPrefix(c.Path(), "/")).Msg("Failed request.")
	}

	return c.Status(code).JSON(CodeResp{Code: code, Message: message})
}

type CodeResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(CodeResp{Code: 200, Message: "Server is up."})
}
