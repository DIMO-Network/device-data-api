package main

import (
	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"fmt"
	"strings"

	"github.com/DIMO-Network/shared/db"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
)

// dependencyContainer way to hold different dependencies we need for our app. We could put all our deps and follow this pattern for everything.
type dependencyContainer struct {
	kafkaProducer sarama.SyncProducer
	settings      *config.Settings
	logger        *zerolog.Logger
	//ddSvc         services.DeviceDefinitionsAPIService
	//deviceSvc     services.DeviceAPIService
	dbs func() *db.ReaderWriter
}

func newDependencyContainer(settings *config.Settings, logger zerolog.Logger, dbs func() *db.ReaderWriter) dependencyContainer {
	return dependencyContainer{
		settings: settings,
		logger:   &logger,
		dbs:      dbs,
	}
}

// getKafkaProducer instantiates a new kafka producer if not already set in our container and returns
func (dc *dependencyContainer) getKafkaProducer() sarama.SyncProducer {
	if dc.kafkaProducer == nil {
		p, err := createKafkaProducer(dc.settings)
		if err != nil {
			dc.logger.Fatal().Err(err).Msg("Could not initialize Kafka producer, terminating")
		}
		dc.kafkaProducer = p
	}
	return dc.kafkaProducer
}

func createKafkaProducer(settings *config.Settings) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = NewJVMCompatiblePartitioner
	p, err := sarama.NewSyncProducer(strings.Split(settings.KafkaBrokers, ","), config)
	if err != nil {
		return nil, fmt.Errorf("failed to construct producer with broker list %s: %w", settings.KafkaBrokers, err)
	}
	return p, nil
}

// getDeviceService dials a grpc connection to devices-api. If fails just returns nil svc, this is ok since not a required svc in all cases
func (dc *dependencyContainer) getDeviceService() (services.DeviceAPIService, *grpc.ClientConn) {
	devicesConn, err := grpc.Dial(dc.settings.DevicesAPIGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		dc.logger.Error().Err(err).Str("devices-api-grpc-addr", dc.settings.DevicesAPIGRPCAddr).
			Msg("failed to dial devices-api grpc - continuing")
		return nil, devicesConn
	}
	deviceSvc := services.NewDeviceAPIService(devicesConn)
	return deviceSvc, devicesConn
}
