package oauth

import (
	"github.com/gofiber/fiber/v2"
)

type CreateOauthRequestBody struct {
	Username string `json:"username" validate:"required,nonempty"`
	Password string `json:"password" validate:"required,nonempty"`
}

func (h *handler) CreateOauth(c *fiber.Ctx) error {
	body := new(CreateOauthRequestBody)

	if err := c.BodyParser(body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	return c.Status(fiber.StatusCreated).JSON("OK")
}
