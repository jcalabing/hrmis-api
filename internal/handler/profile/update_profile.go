package profile

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
)

type CreateRequestBody struct {
	Firstname  string `json:"Firstname" validate:"required,nonempty"`
	Lastname   string `json:"Lastname" validate:"required,nonempty"`
	Middlename string `json:"Middlename" validate:"required,nonempty"`
}

func (h *handler) UpdateProfile(c *fiber.Ctx) error {
	body := new(CreateRequestBody)

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

	fieldlist := []string{"Firstname", "Lastname", "Middlename"}

	profile := c.Locals("profile").(model.User)
	for _, fieldName := range fieldlist {
		fieldValue := reflect.ValueOf(body).Elem().FieldByName(fieldName)
		if fieldValue.IsValid() {
			userfield := model.UserField{
				Key:    fieldName,
				Value:  fieldValue.Interface().(string),
				UserID: profile.ID,
			}
			// Look for User
			var oldUserField model.UserField
			h.DB.Where("user_id = ? AND key = ?", userfield.UserID, userfield.Key).First(&oldUserField)

			// If key does not exist, create a new UserField
			if oldUserField.ID == 0 {
				if result := h.DB.Create(&userfield); result.Error != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while creating the new user field.",
					))
				}
			} else {
				// If key exists, update the existing UserField
				if result := h.DB.Model(&oldUserField).Updates(&userfield); result.Error != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating the user field.",
					))
				}
			}
		}

	}

	return c.Status(fiber.StatusOK).JSON(&profile)

}
