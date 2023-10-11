package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/model"
)

func (h handler) GetUsers(c *fiber.Ctx) error {
	var user []model.User

	if result := h.DB.Find(&user); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	return c.Status(fiber.StatusOK).JSON(&user)

}
