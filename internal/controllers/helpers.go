package controllers

import (
	"github.com/DIMO-Network/device-definitions-api/pkg/grpc"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/volatiletech/null/v8"
)

type QueryValues struct {
	UserID       string `query:"-" json:"userId"`
	UserDeviceID string `query:"-" json:"userDeviceId"`
}

func GetUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}

// duplicate code with devices-api

type DeviceAttributeType string

const (
	MPG                 DeviceAttributeType = "mpg"
	FuelTankCapacityGal DeviceAttributeType = "fuel_tank_capacity_gal"
	MpgHighway          DeviceAttributeType = "mpg_highway"
)

type DeviceDefinitionRange struct {
	FuelTankCapGal float64 `json:"fuel_tank_capacity_gal"`
	Mpg            float64 `json:"mpg"`
	MpgHwy         float64 `json:"mpg_highway"`
}

func GetActualDeviceDefinitionMetadataValues(dd *grpc.GetDeviceDefinitionItemResponse, deviceStyleID null.String) *DeviceDefinitionRange {

	var fuelTankCapGal, mpg, mpgHwy float64 = 0, 0, 0

	var metadata []*grpc.DeviceTypeAttribute

	if !deviceStyleID.IsZero() {
		for _, style := range dd.DeviceStyles {
			if style.Id == deviceStyleID.String {
				metadata = style.DeviceAttributes
				break
			}
		}
	}

	if len(metadata) == 0 && dd != nil && dd.DeviceAttributes != nil {
		metadata = dd.DeviceAttributes
	}

	for _, attr := range metadata {
		switch DeviceAttributeType(attr.Name) {
		case FuelTankCapacityGal:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				fuelTankCapGal = v
			}
		case MPG:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				mpg = v
			}
		case MpgHighway:
			if v, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				mpgHwy = v
			}
		}
	}

	return &DeviceDefinitionRange{
		FuelTankCapGal: fuelTankCapGal,
		Mpg:            mpg,
		MpgHwy:         mpgHwy,
	}
}

func CalculateRange(rangeData *DeviceDefinitionRange, fuelPercentRemaining float64) *float64 {
	// todo refactor this into a pkg library we can share with devices-api so that all projects use same logic
	if rangeData.FuelTankCapGal > 0 && rangeData.Mpg > 0 {
		fuelTankAtGal := rangeData.FuelTankCapGal * fuelPercentRemaining // future: subtract a little for more conservative estimate
		rangeMiles := rangeData.Mpg * fuelTankAtGal
		rangeKm := 1.60934 * rangeMiles
		return &rangeKm
	}
	return nil
}
