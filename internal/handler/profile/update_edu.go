package profile

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
	"gorm.io/gorm"
)

type edu struct {
	Award     string      `json:"award"`
	Degree    string      `json:"degree"`
	Ended     string      `json:"ended"`
	From      string      `json:"from"`
	Graduated string      `json:"graduated"`
	Highest   string      `json:"highest"`
	ID        interface{} `json:"id"`
	Level     string      `json:"level"`
	School    string      `json:"school"`
}

func UpdateEdu(h *handler, c *fiber.Ctx, profile model.User, education string) error {
	var eduArray []edu

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(education), &eduArray); err != nil {
		fmt.Println("Error unmarshaling Education:", err)
	}

	fieldlist := []string{
		"Award",
		"Degree",
		"Ended",
		"From",
		"Graduated",
		"Highest",
		"ID",
		"Level",
		"School",
	}

	var retainEdu []interface{}
	for _, eduData := range eduArray {
		// create an education profile
		eduAttain := model.Edu{
			UserID: profile.ID,
		}
		if eduData.ID == "" {
			if result := h.DB.Create(&eduAttain); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new edu field.",
				))
			}
		} else {
			if result := h.DB.First(&eduAttain, eduData.ID); result.Error != nil {
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Education Not Found",
				))
			}
		}
		for _, fieldName := range fieldlist {
			fieldValue := reflect.ValueOf(eduData).FieldByName(fieldName)
			if fieldValue.IsValid() {
				edufield := model.EduField{
					Key:   fieldName,
					Value: fmt.Sprintf("%v", fieldValue.Interface()),
					EduID: eduAttain.ID,
				}
				// Look for EduField
				var oldEduField model.EduField
				if err := h.DB.Where("edu_id = ? AND key = ?", edufield.EduID, edufield.Key).First(&oldEduField).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// If key does not exist, create a new EduField
						if err := h.DB.Create(&edufield).Error; err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
								fiber.StatusInternalServerError,
								"",
								"Error occurred while creating the new edu field.",
							))
						}
					} else {
						return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
							fiber.StatusInternalServerError,
							"",
							"Error occurred while querying the database.",
						))
					}
				} else {
					// If key exists, update the existing EduField
					if err := h.DB.Model(&oldEduField).Updates(&edufield).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
							fiber.StatusInternalServerError,
							"",
							"Error occurred while updating the edu field.",
						))
					}
				}
			}
		}

		retainEdu = append(retainEdu, eduAttain.ID)
	}

	if err := h.DB.Preload("Educations").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {

		for _, edu := range profile.Educations {
			if !functions.ArrayContains(retainEdu, edu.ID) {
				if err := h.DB.Delete(&edu).Error; err != nil {
					fmt.Println("Error Deleting Education", err)
				}
			}
		}

	}

	return nil
}
