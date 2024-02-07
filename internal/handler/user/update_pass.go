package user

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePassRequestBody struct {
	Username    string `json:"username" validate:"required,nonempty" gorm:"unique"`
	Email       string `json:"email" validate:"required,email,nonempty" gorm:"unique"`
	Password    string `json:"password" validate:"required,nonempty"`
	Newpassword string `json:"newpassword"`
	Repassword  string `json:"repassword"`
	Active      string `json:"active" validate:"required,nonempty"`
}

func (h *handler) UpdatePass(c *fiber.Ctx) error {
	body := new(UpdatePassRequestBody)

	fmt.Println(body)
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
			fmt.Println(validationErrors)
		}
		errorMsg := "Kindly check the following fields: " + strings.Join(validationErrors, ", ")

		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			errorMsg,
			validationErrors,
		))
	}

	var user model.User
	if result := h.DB.First(&user, c.Params("id")); result.Error != nil {
		// return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"User Not Found",
		))
	}

	//verify user password
	profile := c.Locals("profile").(model.User)
	err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(body.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())
	}

	if body.Newpassword != "" {
		if len(body.Newpassword) >= 1 && len(body.Newpassword) <= 5 {
			return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
				fiber.StatusBadRequest,
				"",
				"Password must be more than 6 characters long.",
			))
		}
		if body.Newpassword != body.Repassword {
			return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
				fiber.StatusBadRequest,
				"",
				"Password is not equal.",
			))
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Newpassword), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
				fiber.StatusInternalServerError,
				"",
				"Failed to hash Password",
			))
		}
		user.Password = string(hash)
	}

	user.Username = body.Username
	user.Email = body.Email
	user.Active = body.Active
	h.DB.Save(&user)

	return c.Status(fiber.StatusCreated).JSON(user)
}
