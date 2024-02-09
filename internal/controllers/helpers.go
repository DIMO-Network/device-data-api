package controllers

import (
	"github.com/DIMO-Network/device-data-api/internal/services"
	"github.com/DIMO-Network/shared/middleware/privilegetoken"
	"github.com/DIMO-Network/shared/privileges"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

// CalculateRange calculates the range given device def MPG info and the fuelPercentRemaining - may be duplicated in devices-api
func CalculateRange(rangeData *services.DeviceDefinitionRange, fuelPercentRemaining float64) *float64 {
	// todo refactor this into a pkg library we can share with devices-api so that all projects use same logic
	if rangeData.FuelTankCapGal > 0 && rangeData.Mpg > 0 {
		fuelTankAtGal := rangeData.FuelTankCapGal * fuelPercentRemaining // future: subtract a little for more conservative estimate
		rangeMiles := rangeData.Mpg * fuelTankAtGal
		rangeKm := 1.60934 * rangeMiles
		return &rangeKm
	}
	return nil
}

// getPrivileges takes a fiber context and returns a slice of prvilieges from the jwt token if they exist.
func getPrivileges(c *fiber.Ctx) []privileges.Privilege {
	claims, ok := c.Locals("tokenClaims").(privilegetoken.CustomClaims)
	if !ok {
		return nil
	}
	privs := make([]privileges.Privilege, len(claims.PrivilegeIDs))
	for i, id := range claims.PrivilegeIDs {
		privs[i] = privileges.Privilege(id)
	}
	return privs
}
