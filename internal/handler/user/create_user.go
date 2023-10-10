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
	// Username string `json:"username" validate:"required,nonempty"`
	Email                string `json:"email" validate:"required,email,nonempty"`
	Password             string `json:"password" validate:"required,nonempty"`
	Active               string `json:"active"`
	Agencyemployeenumber string `json:"agencyemployeenumber"`
	Assignment           string `json:"assignment"`
	Civilstatus          string `json:"civilstatus"`
	Dailyrate            string `json:"dailyrate"`
	Dateofbirth          string `json:"dateofbirth"`
	Division             string `json:"division"`
	Extname              string `json:"extname"`
	Firstname            string `json:"firstname"`
	Itemnumber           string `json:"itemnumber"`
	Middlename           string `json:"middlename"`
	Office               string `json:"office"`
	Plantilla            string `json:"plantilla"`
	Positiontitle        string `json:"positiontitle"`
	Remarks              string `json:"remarks"`
	Salarygrade          string `json:"salarygrade"`
	Sex                  string `json:"sex"`
	Sgstep               string `json:"sgstep"`
	Surname              string `json:"surname"`
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
		// Username: body.Username,
		Username: body.Email,
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
