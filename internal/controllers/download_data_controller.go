package controllers

import (
	"fmt"

	"github.com/DIMO-Network/device-data-api/internal/config"
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// DataDownloadController provides endpoints for user to download their data or save it (encrypted) to IPFS
type DataDownloadController struct {
	log         *zerolog.Logger
	querySvc    *services.UserDataService
	ipfsAddress string
	deviceAPI   services.DeviceAPIService
}

func NewDataDownloadController(settings *config.Settings, log *zerolog.Logger, querySvc *services.UserDataService, deviceAPIService services.DeviceAPIService) *DataDownloadController {
	return &DataDownloadController{ipfsAddress: settings.IPFSAddress, log: log, querySvc: querySvc, deviceAPI: deviceAPIService}
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
	err = d.querySvc.UserDataJSONS3(userDeviceID, params.EncryptionKey, params.RangeStart, params.RangeEnd, d.ipfsAddress, params.IPFS)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(map[string]string{"success": "data can be downloaded via links sent to user email on file"})
}
