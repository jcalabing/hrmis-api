package user

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
		return fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body")
	}

	validate := validator.New()
	validate.RegisterValidation("nonempty", validateNonEmpty)

	if err := validate.Struct(body); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field())
		}
		errorMsg := "The following fields are required: " + strings.Join(validationErrors, ", ")
		return fiber.NewError(fiber.StatusBadRequest, errorMsg)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash Password")
	}

	user := model.User{
		Username: body.Username,
		Password: string(hash),
		Email:    body.Email,
	}

	if result := h.DB.Create(&user); result.Error != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create User")
	}

	return c.Status(fiber.StatusCreated).JSON(&user)
}

func validateNonEmpty(fl validator.FieldLevel) bool {
	field := fl.Field()
	value := strings.TrimSpace(field.String())
	return value != ""
}
