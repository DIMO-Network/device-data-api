package fingerprint

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/helpers"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/device-data-api/models"
	"github.com/DIMO-Network/shared"
	"github.com/DIMO-Network/shared/db"
	"github.com/DIMO-Network/shared/kafka"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"regexp"
	"strings"
)

type Event struct {
	shared.CloudEvent[json.RawMessage]
	Signature string `json:"signature"`
}

type Consumer struct {
	logger           *zerolog.Logger
	DBS              db.Store
	deviceAPIService services.DeviceAPIService
}

func NewConsumer(dbs db.Store, log *zerolog.Logger, deviceAPIService services.DeviceAPIService) *Consumer {
	return &Consumer{
		DBS:              dbs,
		logger:           log,
		deviceAPIService: deviceAPIService,
	}
}

func RunConsumer(ctx context.Context, settings *config.Settings, logger *zerolog.Logger, dbs db.Store, deviceAPIService services.DeviceAPIService) error {
	consumer := NewConsumer(dbs, logger, deviceAPIService)

	if err := kafka.Consume(ctx, kafka.Config{
		Brokers: strings.Split(settings.KafkaBrokers, ","),
		Topic:   settings.DeviceFingerprintTopic,
		Group:   settings.DeviceFingerprintConsumerGroup,
	}, consumer.HandleDeviceFingerprint, logger); err != nil {
		logger.Fatal().Err(err).Msg("couldn't start device fingerprint consumer")
	}

	logger.Info().Msg("Starting fingerprint consumer to get VIN.")

	return nil
}

var userDeviceDataPrimaryKeyColumns = []string{models.UserDeviceDatumColumns.UserDeviceID, models.UserDeviceDatumColumns.IntegrationID}

func (c *Consumer) HandleDeviceFingerprint(ctx context.Context, event *Event) error {
	if !common.IsHexAddress(event.Subject) {
		return fmt.Errorf("subject %q not a valid address", event.Subject)
	}
	addr := common.HexToAddress(event.Subject)
	signature := common.FromHex(event.Signature)
	hash := crypto.Keccak256Hash(event.Data)

	if recAddr, err := helpers.Ecrecover(hash.Bytes(), signature); err != nil {
		return fmt.Errorf("failed to recover an address: %w", err)
	} else if recAddr != addr {
		return fmt.Errorf("recovered wrong address %s", recAddr)
	}

	vin, err := ExtractVIN(event.Data)
	if err != nil {
		if err == ErrNoVIN {
			return nil
		}
		return fmt.Errorf("couldn't extract vin: %w", err)
	}

	ud, err := c.deviceAPIService.GetUserDeviceByEthAddr(ctx, addr.Bytes())

	if err != nil {
		return fmt.Errorf("failed querying for device: %w addr %s", err, addr)
	}

	udd, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(ud.Id),
	).One(ctx, c.DBS.DBS().Reader)

	if err != nil {
		return fmt.Errorf("failed querying for device data: %w", err)
	}

	// extract signals with timestamps and persist to signals
	signals := make(map[string]any)
	if err := udd.Signals.Unmarshal(&signals); err != nil {
		return err
	}

	if vinData, ok := signals["vin"].(map[string]interface{}); ok {
		if _, exists := vinData["value"]; exists {
			vinData["value"] = vin
		}
	}

	if vinData, ok := signals["vin"].(map[string]interface{}); ok {
		if _, exists := vinData["timestamp"]; exists {
			vinData["timestamp"] = vin
		}
	}

	if err := udd.Signals.Marshal(signals); err != nil {
		return err
	}

	if err := udd.Upsert(ctx, c.DBS.DBS().Writer, true, userDeviceDataPrimaryKeyColumns, boil.Infer(), boil.Infer()); err != nil {
		return fmt.Errorf("error upserting datum: %w", err)
	}

	return nil
}

var ErrNoVIN = errors.New("no VIN field")
var basicVINExp = regexp.MustCompile(`^[A-Z0-9]{17}$`)

// ExtractVIN extracts the vin field from a status update's data object.
// If this field is not present or fails basic validation, an error is returned.
// The function does clean up the input slightly.
func ExtractVIN(data []byte) (string, error) {
	partialData := new(struct {
		VIN *string `json:"vin"`
	})

	if err := json.Unmarshal(data, partialData); err != nil {
		return "", fmt.Errorf("failed parsing data field: %w", err)
	}

	if partialData.VIN == nil {
		return "", ErrNoVIN
	}

	// Minor cleaning.
	vin := strings.ToUpper(strings.ReplaceAll(*partialData.VIN, " ", ""))

	// We have seen crazy VINs like "\u000" before.
	if !basicVINExp.MatchString(vin) {
		return "", errors.New("invalid VIN")
	}

	return vin, nil
}
