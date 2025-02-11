package controllers

import (
	"time"

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

// parseDateRange validates whether the start and end date parameters are in an acceptable format, otherwise returns an error.
// accepted formats include time.RFC3339 and time.DateOnly
func parseDateRange(startDate, endDate string) (string, string, error) {
	if startDate == "" {
		startDate = time.Now().Add(-1 * (time.Hour * 24 * 14)).Format(time.RFC3339) // if no startdate default to 2 weeks
	} else {
		if !validDate(startDate) {
			return "", "", fiber.NewError(fiber.StatusBadRequest, "Invalid start date format")
		}
	}

	if endDate == "" {
		endDate = time.Now().Format(time.RFC3339)
	} else {
		if !validDate(endDate) {
			return "", "", fiber.NewError(fiber.StatusBadRequest, "Invalid end date format")
		}
	}

	return startDate, endDate, nil
}

// validDate returns true if date can be parsed as time.DateOnly or time.RFC3339
func validDate(date string) bool {
	_, err := time.Parse(time.DateOnly, date)
	if err != nil {
		_, err = time.Parse(time.RFC3339, date)
		if err != nil {
			return false
		}
	}
	return true
}
