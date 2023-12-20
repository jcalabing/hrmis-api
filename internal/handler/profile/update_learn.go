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

type learn struct {
	Title     string      `json:"title"`
	Started   string      `json:"started"`
	Ended     string      `json:"ended"`
	Conducted string      `json:"conducted"`
	Type      string      `json:"type"`
	Hours     string      `json:"hours"`
	ID        interface{} `json:"id"`
}

func UpdateLearn(h *handler, c *fiber.Ctx, profile model.User, learns string) error {
	var learnArray []learn

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(learns), &learnArray); err != nil {
		fmt.Println("Error unmarshaling Learn:", err)
	}

	fieldlist := []string{
		"Title",
		"Started",
		"Ended",
		"Conducted",
		"Type",
		"Hours",
		"ID",
	}

	var retainLearn []interface{}
	for _, learnData := range learnArray {
		// create an learns profile
		learnAttain := model.Learn{
			UserID: profile.ID,
		}
		if learnData.ID == "" {
			if result := h.DB.Create(&learnAttain); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new learn field.",
				))
			}
		} else {
			if result := h.DB.First(&learnAttain, learnData.ID); result.Error != nil {
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Learn Not Found",
				))
			}
		}
		for _, fieldName := range fieldlist {
			fieldValue := reflect.ValueOf(learnData).FieldByName(fieldName)
			if fieldValue.IsValid() {
				learnfield := model.LearnField{
					Key:     fieldName,
					Value:   fmt.Sprintf("%v", fieldValue.Interface()),
					LearnID: learnAttain.ID,
				}
				// Look for LearnField
				var oldLearnField model.LearnField
				if err := h.DB.Where("learn_id = ? AND key = ?", learnfield.LearnID, learnfield.Key).First(&oldLearnField).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// If key does not exist, create a new LearnField
						if err := h.DB.Create(&learnfield).Error; err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
								fiber.StatusInternalServerError,
								"",
								"Error occurred while creating the new learn field.",
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
					// If key exists, update the existing LearnField
					if err := h.DB.Model(&oldLearnField).Updates(&learnfield).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
							fiber.StatusInternalServerError,
							"",
							"Error occurred while updating the learn field.",
						))
					}
				}
			}
		}

		retainLearn = append(retainLearn, learnAttain.ID)
	}

	if err := h.DB.Preload("Learns").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {

		for _, learn := range profile.Learns {
			if !functions.ArrayContains(retainLearn, learn.ID) {
				if err := h.DB.Delete(&learn).Error; err != nil {
					fmt.Println("Error Deleting Learn", err)
				}
			}
		}

	}

	return nil
}
