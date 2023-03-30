package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type QueryValues struct {
	RangeStart    string `query:"rangestart"`
	RangeEnd      string `query:"rangeend"`
	Timezone      string `query:"timezone"`
	EncryptionKey string `json:"encryptionkey"`
	IPFS          bool   `json:"ipfs"`
	UserID        string
	UserDeviceID  string
}

func ValidateQueryParams(p *QueryValues, c *fiber.Ctx) error {
	err := c.QueryParser(p)
	if err != nil {
		return err
	}

	// if empty range start, default to 1/1/2022
	if p.RangeStart == "" {
		p.RangeStart = "20220101"
	}

	rangeStart, err := time.Parse("20060102", p.RangeStart)
	if err != nil {
		return err
	}
	p.RangeStart = rangeStart.Format("2006-01-02")

	// if empty range end, default to current time and return
	if p.RangeEnd == "" {
		p.RangeEnd = "now"
		return nil
	}

	// if passed rangeend, parse and return
	rangeEnd, err := time.Parse("20060102", p.RangeEnd)
	if err != nil {
		return err
	}
	p.RangeEnd = rangeEnd.Format("2006-01-02")

	return nil
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}
