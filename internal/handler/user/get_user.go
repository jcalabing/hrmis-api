package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/model"
)

func (h handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var user model.User

	if result := h.DB.First(&user, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	returnValue := map[string]any{
		"user": model.ConvertToUserResponse(h.DB, user),
	}

	return c.Status(fiber.StatusOK).JSON(&returnValue)
	// return c.Status(fiber.StatusOK).JSON(&user)

}
