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
	querySvc    *services.AggregateQueryService
	ipfsAddress string
	deviceAPI   services.DeviceAPIService
}

func NewDataDownloadController(settings *config.Settings, log *zerolog.Logger, querySvc *services.AggregateQueryService, deviceAPIService services.DeviceAPIService) *DataDownloadController {
	return &DataDownloadController{ipfsAddress: settings.IPFSAddress, log: log, querySvc: querySvc, deviceAPI: deviceAPIService}
}

// DownloadHandler godoc
// @Description  returns all user data
// @Tags         devices
// @Produce      json
// @Success      200  {object}  []models.UserData
// @Router       /download [get]
func (d *DataDownloadController) DownloadHandler(c *fiber.Ctx) error {
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
	response, err := d.querySvc.DownloadUserData(userDeviceID, params.EncryptionKey, params.RangeStart, params.RangeEnd, d.ipfsAddress, params.IPFS)
	if err != nil {
		return c.JSON(map[string]string{"Error": err.Error()})
	}
	return c.JSON(&response)
}
