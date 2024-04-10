package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/DIMO-Network/device-data-api/internal/services/elastic"
)

// historicalV2QueryParams represents the query parameters for the historical endpoint.
type historicalV2QueryParams struct {
	StartTime time.Time `query:"startTime"`
	EndTime   time.Time `query:"endTime"`

	// uint32 is used to prevent negative values
	Buckets uint32 `query:"buckets"`
}

type historyResp struct {
	Statuses []json.RawMessage `json:"statuses"`
}

// GetHistoricalPermissioned godoc
// @Description  Get all historical data for a tokenID, within start and end range
// @Tags         device-data
// @Produce      json
// @Success      200
// @Param        tokenID  path   int64  true   "token id"
// @Param        startTime     query  string  false  "startTime is an RFC3339 formatted date-time string. If empty two weeks from endTime"
// @Param        endTime       query  string  false  "endTime is an RFC3339 formatted date-time string. If empty the current time"
// @Param        buckets       query  string  false  "number of data points to return, default 1000"
// @Security     BearerAuth
// @Router       /v2/vehicle/{tokenID}/history [get]
func (d *DeviceDataControllerV2) GetHistoricalPermissioned(c *fiber.Ctx) error {
	// get the startTime and endTime from the query parameters and convert them to time.Time
	var err error
	qParams := historicalV2QueryParams{}
	if err = c.QueryParser(&qParams); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("error parsing query parameters: %v", err))
	}

	// get vehicleID from tokenID
	tokenIDParam := c.Params("tokenID")
	tokenID, err := strconv.ParseInt(tokenIDParam, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("error parsing tokenID: %v", err))
	}
	ctx, cancel := context.WithTimeout(c.Context(), defaultTimeout)
	defer cancel()
	userDevice, err := d.deviceAPI.GetUserDeviceByTokenID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("failed to get user device by token id: %w", err)
	}

	// get the token claims from the context
	privs := getPrivileges(c)
	if len(privs) == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "")
	}

	params := elastic.GetHistoryParams{
		DeviceID:     userDevice.GetId(),
		StartTime:    qParams.StartTime,
		EndTime:      qParams.EndTime,
		Buckets:      int(qParams.Buckets),
		PrivilegeIDs: privs,
	}

	params.SetDefaultHistoryParams()

	result, err := d.esService.GetHistory(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}
	resp := historyResp{Statuses: result}

	// TODO: must be updated to work with modified query
	// body, err =  IfNotExists(c.Context(), d.definitionsAPI, body, userDevice.DeviceDefinitionId, userDevice.DeviceStyleId)
	// if err != nil {
	// 	localLog.Warn().Err(err).Msg("could not add range calculation to document")
	// }
	// body = removeOdometerIfInvalid(body)

	err = c.Status(fiber.StatusOK).JSON(resp)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	return nil
}
