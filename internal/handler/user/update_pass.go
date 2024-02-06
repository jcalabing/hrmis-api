package user

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePassRequestBody struct {
	Email    string `json:"email" validate:"required,email,nonempty"`
	Password string `json:"password"`
	Active   string `json:"active"`
}

func (h *handler) UpdatePass(c *fiber.Ctx) error {
	body := new(UpdatePassRequestBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"",
		))
	}

	validate := validator.New()
	validate.RegisterValidation("nonempty", functions.ValidateNonEmpty)

	if err := validate.Struct(body); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field())
		}
		errorMsg := "Kindly check the following fields: " + strings.Join(validationErrors, ", ")

		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			errorMsg,
			validationErrors,
		))
	}

	var profile model.User
	if c.Params("id") != "" {
		if result := h.DB.First(&profile, c.Params("id")); result.Error != nil {
			// return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
			return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
				fiber.StatusBadRequest,
				"",
				"User Not Found",
			))
		}
	} else {
		profile = c.Locals("profile").(model.User)
	}

	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
				fiber.StatusInternalServerError,
				"",
				"Failed to hash Password",
			))
		}
		profile.Password = string(hash)
	}

	profile.Username = body.Email
	profile.Email = body.Email
	profile.Active = body.Active
	h.DB.Save(&profile)

	return c.Status(fiber.StatusCreated).JSON(profile)
}
