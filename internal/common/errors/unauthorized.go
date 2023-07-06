package errors

import (
	"github.com/gofiber/fiber/v2"
)

func UnAuthorized() *ErrorResponse {
	return NewErrorResponse(
		fiber.StatusUnauthorized,
		"",
		"Invalid Credentials",
	)
}
