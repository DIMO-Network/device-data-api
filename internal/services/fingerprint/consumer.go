package fingerprint

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

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
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Event struct {
	shared.CloudEvent[json.RawMessage]
	Signature string `json:"signature"`
}

type Consumer struct {
	logger           *zerolog.Logger
	dbs              db.Store
	deviceAPIService services.DeviceAPIService
}

func NewConsumer(dbs db.Store, log *zerolog.Logger, deviceAPIService services.DeviceAPIService) *Consumer {
	return &Consumer{
		dbs:              dbs,
		logger:           log,
		deviceAPIService: deviceAPIService,
	}
}

func RunConsumer(ctx context.Context, settings *config.Settings, logger *zerolog.Logger, dbs db.Store, deviceAPIService services.DeviceAPIService) error {
	consumer := NewConsumer(dbs, logger, deviceAPIService)

	if err := kafka.Consume[*Event](ctx, kafka.Config{
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

const autoPiIntegrationID = "27qftVRWQYpVDcO5DltO5Ojbjxk"

// HandleDeviceFingerprint extracts the VIN from the payload, and sets it in user_device_data.signals
func (c *Consumer) HandleDeviceFingerprint(ctx context.Context, event *Event) error {
	if !common.IsHexAddress(event.Subject) {
		return fmt.Errorf("subject %q not a valid eth hex address", event.Subject)
	}
	addr := common.HexToAddress(event.Subject)
	signature := common.FromHex(event.Signature)
	hash := crypto.Keccak256Hash(event.Data)

	if recAddr, err := helpers.Ecrecover(hash.Bytes(), signature); err != nil {
		return fmt.Errorf("fingerprint failed to recover an address: %w", err)
	} else if recAddr != addr {
		return fmt.Errorf("fingerprint recovered wrong address: %s. subject addr: %s", recAddr, addr)
	}

	vin, err := ExtractVIN(event.Data)
	if err != nil {
		if errors.Is(err, ErrNoVIN) {
			return nil
		}
		return fmt.Errorf("fingerprint couldn't extract vin: %w. subject addr: %s", err, addr)
	}

	ud, err := c.deviceAPIService.GetUserDeviceByEthAddr(ctx, addr.Bytes())

	if err != nil {
		return fmt.Errorf("failed querying for device: %w addr %s", err, addr)
	}

	tx, err := c.dbs.DBS().Writer.BeginTx(ctx, nil)
	defer tx.Rollback() // nolint
	if err != nil {
		return errors.Wrap(err, "unable to start transaction")
	}

	udd, err := models.UserDeviceData(
		models.UserDeviceDatumWhere.UserDeviceID.EQ(ud.Id),
		models.UserDeviceDatumWhere.IntegrationID.EQ(autoPiIntegrationID),
	).One(ctx, tx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// create new
			udd = &models.UserDeviceDatum{
				UserDeviceID:  ud.Id,
				IntegrationID: autoPiIntegrationID,
				Signals:       null.JSONFrom([]byte(`{}`)),
			}
		} else {
			return fmt.Errorf("failed querying for device data: %w", err)
		}
	}
	// set vin value in json
	j, err := sjson.SetBytes(udd.Signals.JSON, "vin.value", vin)
	if err != nil {
		return errors.Wrap(err, "failed to set vin value in signals json")
	}
	// set vin timestamps in json
	ts := gjson.GetBytes(event.Data, "timestamp").Int()
	if ts > 0 {
		j, err = sjson.SetBytes(j, "vin.timestamp", time.UnixMilli(ts).UTC().Format(time.RFC3339))
		if err != nil {
			return errors.Wrap(err, "failed to set timestamp value in signals json")
		}
	}
	udd.Signals = null.JSONFrom(j)

	err = udd.Upsert(ctx, tx, true, userDeviceDataPrimaryKeyColumns, boil.Infer(), boil.Infer())
	if err != nil {
		return errors.Wrap(err, "failed to upsert device data for vin")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "failed to commit trx for device data")
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
