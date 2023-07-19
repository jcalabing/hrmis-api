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

type CreateUserRequestBody struct {
	Username string `json:"username" validate:"required,nonempty"`
	Email    string `json:"email" validate:"required,email,nonempty"`
	Password string `json:"password" validate:"required,nonempty"`
}

func (h *handler) CreateUser(c *fiber.Ctx) error {
	body := new(CreateUserRequestBody)

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

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"",
			"Failed to hash Password",
		))
	}

	user := model.User{
		Username: body.Username,
		Password: string(hash),
		Email:    body.Email,
	}

	if result := h.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"",
			"The username or email are not available.",
		))
	}

	return c.Status(fiber.StatusCreated).JSON("OK")
}

// func validateNonEmpty(fl validator.FieldLevel) bool {
// 	field := fl.Field()
// 	value := strings.TrimSpace(field.String())
// 	return value != ""
// }
