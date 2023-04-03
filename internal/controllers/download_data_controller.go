package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// DataDownloadController provides endpoints for user to download their data or save it (encrypted) to IPFS
type DataDownloadController struct {
	log        *zerolog.Logger
	QuerySvc   *services.DataQueryService
	StorageSvc *services.StorageService
	EmailSvc   *services.EmailService
	NATSSvc    *services.NATSService
	deviceAPI  services.DeviceAPIService
}

func NewDataDownloadController(settings *config.Settings, log *zerolog.Logger, esClient8 *es8.TypedClient, deviceAPIService services.DeviceAPIService) (*DataDownloadController, error) {
	querySvc := services.NewAggregateQueryService(esClient8, settings, log)
	storageSvc, err := services.NewStorageService(settings, log)
	if err != nil {
		return nil, err
	}
	emailSvc := services.NewEmailService(settings, log)
	nats, err := services.NewNATSService(settings, log)
	if err != nil {
		return nil, err
	}
	return &DataDownloadController{
		log:        log,
		QuerySvc:   querySvc,
		StorageSvc: storageSvc,
		EmailSvc:   emailSvc,
		deviceAPI:  deviceAPIService,
		NATSSvc:    nats}, nil
}

// DataDownloadHandler godoc
// @Description  returns message indicating that download will be sent to user email
// @Tags         device-data
// @Produce      json
// @Success      200
// @Router       /user/device-data/:userDeviceID/export/json/email [get]
func (d *DataDownloadController) DataDownloadHandler(c *fiber.Ctx) error {
	userID := getUserID(c)
	userDeviceID := c.Params("userDeviceID")
	exists, err := d.deviceAPI.UserDeviceBelongsToUserID(c.Context(), userID, userDeviceID)
	if err != nil {
		return err
	}
	if !exists {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("No device %s found for user %s", userDeviceID, userID))
	}

	var params QueryValues
	err = ValidateQueryParams(&params, c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	params.UserDeviceID = userDeviceID
	params.UserID = userID

	b, err := json.Marshal(params)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = d.NATSSvc.JetStream.Publish(d.QuerySvc.Settings.NATSDataDownloadSubject, b)
	if err != nil {
		return err
	}

	return c.JSON(dataDownloadRequestStatus{
		Status:       "success",
		UserID:       userID,
		UserDeviceID: userDeviceID,
		Message:      "your request has been received; data will be sent to the email associated with your account",
	})
}

func (d *DataDownloadController) DataDownloadConsumer(ctx context.Context) error {
	sub, err := d.NATSSvc.JetStream.PullSubscribe(d.NATSSvc.JetStreamSubject, d.NATSSvc.DurableConsumer, nats.AckWait(d.NATSSvc.AckTimeout))
	if err != nil {
		return err
	}

	for {
		msgs, err := sub.Fetch(1)
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			d.log.Err(err).Msg("error fetching from data download stream")
		}

		for _, msg := range msgs {
			mtd, err := msg.Metadata()
			if err != nil {
				d.log.Info().Err(err).Msg("unable to parse metadata for message")
			}

			select {
			case <-ctx.Done():
				if err := msg.Nak(); err != nil {
					d.log.Info().Err(err).Msgf("data download msg.Nak failure")
				}
				return nil
			default:
				var params QueryValues
				err = json.Unmarshal(msg.Data, &params)
				if err != nil {
					if err := msg.Nak(); err != nil {
						d.log.Error().Msgf("message nak failed: %+v", err)
						return err
					}
					d.log.Error().Msgf("unable to parse query parameters: %+v", err)
					continue
				}

				d.log.Info().Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("data download initiated")
				msg.InProgress()

				nestedCtx, cancel := context.WithCancel(ctx)
				defer cancel()

				go func() {

					tick := time.NewTicker(1 * time.Second)
					defer tick.Stop()
					for {
						select {
						case <-nestedCtx.Done():
							return
						case <-tick.C:
							msg.InProgress()
						}
					}
				}()

				ud, err := d.QuerySvc.FetchUserData(params.UserDeviceID, params.RangeStart, params.RangeEnd, params.Timezone)
				if err != nil {
					d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("error while fetching data from elasticsearch")
					cancel()

					if err := msg.Nak(); err != nil {
						d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("error while calling Nak")
					}

					continue
				}
				cancel()

				nestedCtx, cancel = context.WithCancel(ctx)
				defer cancel()

				go func() {

					tick := time.NewTicker(1 * time.Second)
					defer tick.Stop()
					for {
						select {
						case <-nestedCtx.Done():
							return
						case <-tick.C:
							msg.InProgress()
						}
					}
				}()

				keyName := fmt.Sprintf("userDownloads/%+v/DIMODeviceData_%+v_%+v_%+v", params.UserDeviceID, params.UserDeviceID, params.RangeStart, params.RangeEnd)
				s3link, err := d.StorageSvc.UploadUserData(ctx, ud, keyName)
				if err != nil {
					d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("error while uploading data to s3")
					cancel()

					if err := msg.Nak(); err != nil {
						d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("error while calling Nak")
					}

					continue
				}
				cancel()

				msg.InProgress()

				err = d.EmailSvc.SendEmail(params.UserID, s3link)
				if err != nil {
					if err := msg.Nak(); err != nil {
						d.log.Error().Msgf("message nak failed: %+v", err)
						return err
					}
					d.log.Err(err).Msg("unable to put send email")
					continue
				}

				msg.Ack()
				d.log.Info().Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Uint64("numDelivered", mtd.NumDelivered).Msg("data download completed")
			}
		}
	}
}

type dataDownloadRequestStatus struct {
	Status       string `json:"status"`
	UserID       string `json:"userId"`
	UserDeviceID string `json:"userDeviceId"`
	Message      string `json:"message"`
	RangeStart   string `json:"rangeStart,omitempty"`
	RangeEnd     string `json:"rangeEnd,omitempty"`
}
