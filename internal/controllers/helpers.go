package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type QueryValues struct {
	RangeStart   string `query:"rangestart"`
	RangeEnd     string `query:"rangeend"`
	Timezone     string `query:"timezone"`
	UserID       string
	UserDeviceID string
}

func ValidateQueryParams(p *QueryValues, c *fiber.Ctx) error {
	err := c.QueryParser(p)
	if err != nil {
		return err
	}

	return nil
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}
