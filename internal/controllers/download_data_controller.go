package controllers

import (
	"fmt"
	"time"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// DataDownloadController provides endpoints for user to download their data or save it (encrypted) to IPFS
type DataDownloadController struct {
	log        *zerolog.Logger
	QuerySvc   *services.DataQueryService
	StorageSvc *services.StorageService
	EmailSvc   *services.EmailService
	deviceAPI  services.DeviceAPIService
}

func NewDataDownloadController(settings *config.Settings, log *zerolog.Logger, esClient8 *es8.TypedClient, deviceAPIService services.DeviceAPIService) *DataDownloadController {
	querySvc := services.NewAggregateQueryService(esClient8, settings, log)
	storageSvc := services.NewStorageService(settings, log)
	emailSvc := services.NewEmailService(settings, log)
	return &DataDownloadController{log: log, QuerySvc: querySvc, StorageSvc: storageSvc, EmailSvc: emailSvc, deviceAPI: deviceAPIService}
}

// JSONDownloadHandler godoc
// @Description  returns user data as json
// @Tags         device-data
// @Produce      json
// @Success      200  {object}
// @Router       /user/device-data/:userDeviceID/export/json/email [get]
func (d *DataDownloadController) JSONDownloadHandler(c *fiber.Ctx) error {
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
	data, err := d.QuerySvc.FetchUserData(userDeviceID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	keyName := "userDownloads/" + userDeviceID + "/" + time.Now().Format(time.RFC3339) + ".json"
	s3link, err := d.StorageSvc.UploadUserData(data, keyName)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = d.EmailSvc.SendEmail(userDeviceID, s3link)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(map[string]string{"success": "data can be downloaded via links sent to user email on file"})
}
