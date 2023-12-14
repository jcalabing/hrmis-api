package profile

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/model"
	"gorm.io/gorm"
)

type edu struct {
	Award     string `json:"award"`
	Degree    string `json:"degree"`
	Ended     string `json:"ended"`
	From      string `json:"from"`
	Graduated string `json:"graduated"`
	Highest   string `json:"highest"`
	ID        string `json:"id"`
	Level     string `json:"level"`
	School    string `json:"school"`
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

	for _, eduData := range eduArray {
		fmt.Println("The Level is" + eduData.Level)
		fmt.Printf("Struct: %+v\n", eduData)
		fmt.Println("-------------------------------LOOP 1------------------------------------")
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
			fmt.Println("HAS Empty ID")
		} else {
			if result := h.DB.First(&eduAttain, eduData.ID); result.Error != nil {
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Education Not Found",
				))
			}
			fmt.Println("HAS VALID ID")
		}
		fmt.Println("-------------------------------LOOP 1 part 1------------------------------------")

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
					fmt.Println("Has data")
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

		fmt.Println("ID Created: ", eduAttain.ID)
	}

	return nil
}
