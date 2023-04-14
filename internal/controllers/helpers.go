package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type QueryValues struct {
	Start        time.Time `query:"start" json:"start"`
	End          time.Time `query:"end" json:"end"`
	UserID       string    `query:"-" json:"userId"`
	UserDeviceID string    `query:"-" json:"userDeviceId"`
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}
