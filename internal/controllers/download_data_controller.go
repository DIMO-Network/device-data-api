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
	log       *zerolog.Logger
	QuerySvc  *services.QueryStorageService
	EmailSvc  *services.EmailService
	NATSSvc   *services.NATSService
	deviceAPI services.DeviceAPIService
}

func NewDataDownloadController(settings *config.Settings, log *zerolog.Logger, esClient8 *es8.TypedClient, deviceAPIService services.DeviceAPIService) (*DataDownloadController, error) {
	querySvc, err := services.NewQueryStorageService(esClient8, settings, log)
	if err != nil {
		return nil, err
	}
	emailSvc := services.NewEmailService(settings, log)
	nats, err := services.NewNATSService(settings, log)
	if err != nil {
		return nil, err
	}
	return &DataDownloadController{
		log:       log,
		QuerySvc:  querySvc,
		EmailSvc:  emailSvc,
		deviceAPI: deviceAPIService,
		NATSSvc:   nats}, nil
}

// DataDownloadHandler godoc
// @Description  Enqueues a data export job for the specified device. A link to download
// @Description  a large JSON file of signals will be emailed to the address on file for the
// @Description  current user.
// @Tags         device-data
// @Produce      json
// @Success      200
// @Security     BearerAuth
// @Param        userDeviceID  path   string  true   "Device id" Example(2OQjmqUt9dguQbJt1WImuVfje3W)
// @Param        start     query  string  false  "Start timestamp in RFC-3339 format" Example(2023-04-14T07:20:50.52Z)
// @Param        end       query  string  false  "End timestamp in RFC-3339 format" Example(2023-04-14T08:11:33.94Z)
// @Router       /user/device-data/{userDeviceID}/export/json/email [get]
func (d *DataDownloadController) DataDownloadHandler(c *fiber.Ctx) error {
	userID := getUserID(c)
	userDeviceID := c.Params("userDeviceID")
	if isOwner, err := d.deviceAPI.UserDeviceBelongsToUserID(c.Context(), userID, userDeviceID); err != nil {
		return err
	} else if !isOwner {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("No device %s found for user %s", userDeviceID, userID))
	}

	var params QueryValues
	if err := c.QueryParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	params.UserDeviceID = userDeviceID
	params.UserID = userID

	b, err := json.Marshal(params)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if _, err = d.NATSSvc.JetStream.Publish(d.QuerySvc.NATSDataDownloadSubject, b); err != nil {
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
				d.nak(msg, nil)
				d.log.Info().Err(err).Msg("unable to parse metadata for message")
				continue
			}

			select {
			case <-ctx.Done():
				d.nak(msg, nil)
				return nil
			default:
				var params QueryValues
				err = json.Unmarshal(msg.Data, &params)
				if err != nil {
					d.nak(msg, &params)
					d.log.Error().Msgf("unable to parse query parameters: %+v", err)
					continue
				}

				d.log.Info().Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("data download initiated")
				d.inProgress(msg, params)

				nestedCtx, cancel := context.WithCancel(ctx)
				go func() {
					tick := time.NewTicker(5 * time.Second)
					defer tick.Stop()
					for {
						select {
						case <-nestedCtx.Done():
							return
						case <-tick.C:
							d.inProgress(msg, params)
						}
					}
				}()

				s3link, err := d.QuerySvc.StreamDataToS3(ctx, params.UserDeviceID, params.Start, params.End)
				if err != nil {
					d.nak(msg, &params)
					d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("error while fetching data from elasticsearch")
					cancel()
					continue
				}
				cancel()

				d.inProgress(msg, params)

				err = d.EmailSvc.SendEmail(params.UserID, s3link)
				if err != nil {
					d.nak(msg, &params)
					d.log.Err(err).Msg("unable to put send email")
					continue
				}

				d.ack(msg, params)
				d.log.Info().Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Uint64("numDelivered", mtd.NumDelivered).Msg("data download completed")
			}
		}
	}
}

func (d *DataDownloadController) ack(msg *nats.Msg, params QueryValues) {
	if err := msg.Ack(); err != nil {
		d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("message ack failed")
	}
}

func (d *DataDownloadController) inProgress(msg *nats.Msg, params QueryValues) {
	if err := msg.InProgress(); err != nil {
		d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("message in progress failed")
	}
}

func (d *DataDownloadController) nak(msg *nats.Msg, params *QueryValues) {
	err := msg.Nak()
	if params == nil {
		d.log.Err(err).Msg("message nak failed")
	} else {
		d.log.Err(err).Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Msg("message nak failed")
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
