package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type QueryValues struct {
	RangeStart   string `query:"rangestart" json:"rangeStart"`
	RangeEnd     string `query:"rangeend" json:"rangeEnd"`
	Timezone     string `query:"timezone" json:"timeZone"`
	UserID       string `query:"-" json:"userId"`
	UserDeviceID string `query:"-" json:"userDeviceId"`
}

func ValidateQueryParams(p *QueryValues, c *fiber.Ctx) error {
	return c.QueryParser(p)
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}
