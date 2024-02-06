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
)

type UpdateUserRequestBody struct {
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

func (h *handler) UpdateUser(c *fiber.Ctx) error {
	body := new(UpdateUserRequestBody)

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

	// if body.Password != "" {
	// 	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
	// 			fiber.StatusInternalServerError,
	// 			"",
	// 			"Failed to hash Password",
	// 		))
	// 	}
	// 	profile.Password = string(hash)
	// }

	// profile.Username = body.Email
	// profile.Email = body.Email
	// profile.Active = body.Active
	// h.DB.Save(&profile)

	// fmt.Println(user.ID)
	//User create fields
	fieldlist := []string{
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
				UserID: profile.ID,
			}

			if result := h.DB.Create(&userfield); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new user field.",
				))
			}

		}

	}

	return c.Status(fiber.StatusCreated).JSON(profile)
}
