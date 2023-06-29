package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/model"
)

func (h handler) GetProfile(c *fiber.Ctx) error {
	profile, _ := c.Locals("profile").(model.User)

	return c.Status(fiber.StatusOK).JSON(&profile)

}
