package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type QueryValues struct {
	RangeStart   string `query:"range_start" json:"rangeStart"`
	RangeEnd     string `query:"range_end" json:"rangeEnd"`
	Timezone     string `query:"time_zone" json:"timeZone"`
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
