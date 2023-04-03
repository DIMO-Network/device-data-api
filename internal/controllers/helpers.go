package controllers

import (
	"time"

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

	err := c.QueryParser(p)

	if p.Timezone == "" {
		p.Timezone = "America/New_York"
	}

	if p.RangeStart == "" {
		p.RangeStart = "20220101"
	}

	if p.RangeEnd == "" {
		p.RangeEnd = time.Now().Format("20060102")
	}

	s, err := time.Parse("20060102", p.RangeStart)
	if err != nil {
		return err
	}

	e, err := time.Parse("20060102", p.RangeEnd)
	if err != nil {
		return err
	}

	p.RangeStart = s.Format("2006-01-02")
	p.RangeEnd = e.Format("2006-01-02")

	return nil
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}
