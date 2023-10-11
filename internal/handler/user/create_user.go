package user

import (
	"fmt"
	"reflect"
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
		Active:   body.Active,
	}

	if result := h.DB.Create(&user); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"",
			"The username or email are not available.",
		))
	}

	fmt.Println(user.ID)
	//User create fields
	fieldlist := []string{"Active",
		"Agencyemployeenumber",
		"Assignment",
		"Civilstatus",
		"Dailyrate",
		"Dateofbirth",
		"Division",
		"Extname",
		"Firstname",
		"Itemnumber",
		"Middlename",
		"Office",
		"Plantilla",
		"Positiontitle",
		"Remarks",
		"Salarygrade",
		"Sex",
		"Sgstep",
		"Surname"}
	for _, fieldName := range fieldlist {
		fieldValue := reflect.ValueOf(body).Elem().FieldByName(fieldName)
		fmt.Println(fieldValue)
		if fieldValue.IsValid() {
			userfield := model.UserField{
				Key:    strings.ToLower(fieldName),
				Value:  fieldValue.Interface().(string),
				UserID: user.ID,
			}
			// Look for User
			// var oldUserField model.UserField
			// h.DB.Where("user_id = ? AND key = ?", userfield.UserID, userfield.Key).First(&oldUserField)

			// If key does not exist, create a new UserField
			// if oldUserField.ID == 0 {
			if result := h.DB.Create(&userfield); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new user field.",
				))
			}
			// } else {
			// 	// If key exists, update the existing UserField
			// 	if result := h.DB.Model(&oldUserField).Updates(&userfield); result.Error != nil {
			// 		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			// 			fiber.StatusInternalServerError,
			// 			"",
			// 			"Error occurred while updating the user field.",
			// 		))
			// 	}
			// }
		}

	}

	return c.Status(fiber.StatusCreated).JSON(user)
	// return c.Status(fiber.StatusCreated).JSON("OK")
}

// func validateNonEmpty(fl validator.FieldLevel) bool {
// 	field := fl.Field()
// 	value := strings.TrimSpace(field.String())
// 	return value != ""
// }
