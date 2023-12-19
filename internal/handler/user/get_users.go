package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/model"
)

func (h handler) GetUsers(c *fiber.Ctx) error {
	var users []model.User

	if result := h.DB.Find(&users); result.Error != nil {
		return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
	}

	var returnValue []any
	for _, user := range users {
		returnValue = append(returnValue, model.ConvertToUserResponse(h.DB, user))
	}

	return c.Status(fiber.StatusOK).JSON(&returnValue)

}
