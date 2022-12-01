package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type QueryValues struct {
	RangeStart    string `query:"rangestart"`
	RangeEnd      string `query:"rangeend"`
	EncryptionKey string `json:"encryptionkey"`
	IPFS          bool   `json:"ipfs"`
}

func ValidateQueryParams(p *QueryValues, c *fiber.Ctx) error {
	err := c.QueryParser(p)
	if err != nil {
		return err
	}
	p.RangeStart, p.RangeEnd, err = validateDateParams(p.RangeStart, p.RangeEnd)
	if err != nil {
		return err
	}
	return nil
}

func validateDateParams(start, end string) (string, string, error) {
	// defaults to past 24 horus if no time range is specified
	if start == "" || end == "" {
		return "2022-01-01", "now", nil
	}
	sd, err := time.Parse("20060102", start)
	if err != nil {
		return "", "", err
	}
	ed, err := time.Parse("20060102", end)
	if err != nil {
		return "", "", err
	}
	return sd.Format("2006-01-02"), ed.Format("2006-01-02"), nil
}

func getUserID(c *fiber.Ctx) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)
	return userID
}
