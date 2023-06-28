package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/model"
)

func (h handler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var book model.User

	if result := h.DB.First(&book, id); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&book)

}
