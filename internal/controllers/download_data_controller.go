package controllers

import (
	"context"
	"encoding/json"
	"errors"
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
	querySvc  *services.QueryStorageService
	emailSvc  *services.EmailService
	natsSvc   *services.NATSService
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
		querySvc:  querySvc,
		emailSvc:  emailSvc,
		deviceAPI: deviceAPIService,
		natsSvc:   nats}, nil
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
// @Router       /v1/user/device-data/{userDeviceID}/export/json/email [post]
func (d *DataDownloadController) DataDownloadHandler(c *fiber.Ctx) error {
	userID := GetUserID(c)
	userDeviceID := c.Params("userDeviceID")
	// make sure user has an email address before enqueueing job
	_, err := d.emailSvc.GetVerifiedEmailAddress(userID)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(dataDownloadRequestStatus{
			Status:       "error",
			UserID:       userID,
			UserDeviceID: userDeviceID,
			Message:      "Your account does not have a verified email address. Please look for your DIMO verification email, click on the verify link and try here again.",
		})
	}

	params := QueryValues{
		UserID:       userID,
		UserDeviceID: userDeviceID,
	}
	b, _ := json.Marshal(params)

	if _, err := d.natsSvc.JetStream.Publish(d.querySvc.NATSDataDownloadSubject, b); err != nil {
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
	sub, err := d.natsSvc.JetStream.PullSubscribe(d.natsSvc.JetStreamSubject, d.natsSvc.DurableConsumer,
		nats.AckWait(d.natsSvc.AckTimeout))
	if err != nil {
		return err
	}
	// it is not possible to set the MaxDeliver in the subscriber here without some additional setup in the js engine, so therefore all message error conditions
	// should be Acked so that we don't end up in infinite retry scenarios. Assuming the error is transient, if the user tries again it will likely work.
	localLog := d.log.With().Str("handler", "DataDownloadConsumer").Logger()

	for {
		msgs, err := sub.Fetch(1)
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}
			localLog.Err(err).Msg("error fetching from data download stream")
		}

		for _, msg := range msgs {
			mtd, err := msg.Metadata()
			if err != nil {
				localLog.Err(err).Msg("unable to parse metadata for message")
				ack(msg, localLog)
				continue
			}

			select {
			case <-ctx.Done():
				ack(msg, localLog)
				return nil
			default:
				var params QueryValues
				err = json.Unmarshal(msg.Data, &params)
				if err != nil {
					ack(msg, localLog)
					localLog.Error().Msgf("unable to parse query parameters: %+v", err)
					continue
				}
				localLog := localLog.With().Str("userId", params.UserID).Str("userDeviceID", params.UserDeviceID).Logger()

				localLog.Info().Msg("data download initiated")
				inProgress(msg, localLog)

				nestedCtx, cancel := context.WithCancel(ctx)
				go func() {
					tick := time.NewTicker(5 * time.Second)
					defer tick.Stop()
					for {
						select {
						case <-nestedCtx.Done():
							return
						case <-tick.C:
							inProgress(msg, localLog)
						}
					}
				}()

				s3link, err := d.querySvc.StreamDataToS3(ctx, params.UserDeviceID)
				if err != nil {
					ack(msg, localLog)
					localLog.Err(err).Msg("error while fetching data from elasticsearch")
					cancel()
					continue
				}
				cancel()

				inProgress(msg, localLog)

				err = d.emailSvc.SendEmail(params.UserID, s3link)
				if err != nil {
					ack(msg, localLog)
					localLog.Err(err).Msg("unable to put send email")
					continue
				}

				ack(msg, localLog)
				localLog.Info().Uint64("numDelivered", mtd.NumDelivered).Msg("data download completed")
			}
		}
	}
}

// ack does msg.Ack() but if there is an error it uses the passed in localLog, which has extra info, to log the error message
// we could use msg.Nak for certain transient error cases but only if we're able to configure the subscriber to have MaxDeliver attempts not be infinite
func ack(msg *nats.Msg, localLog zerolog.Logger) {
	if err := msg.Ack(); err != nil {
		localLog.Err(err).Msg("message ack failed")
	}
}

func inProgress(msg *nats.Msg, localLog zerolog.Logger) {
	if err := msg.InProgress(); err != nil {
		localLog.Err(err).Msg("message in progress failed")
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
