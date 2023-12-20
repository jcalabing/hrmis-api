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

type vol struct {
	Organization string      `json:"organization"`
	Started      string      `json:"started"`
	Ended        string      `json:"ended"`
	Position     string      `json:"position"`
	Hours        string      `json:"hours"`
	ID           interface{} `json:"id"`
}

func UpdateVol(h *handler, c *fiber.Ctx, profile model.User, voluntaries string) error {
	var volArray []vol

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(voluntaries), &volArray); err != nil {
		fmt.Println("Error unmarshaling Voluntaries:", err)
	}

	fieldlist := []string{
		"Organization",
		"Started",
		"Ended",
		"Position",
		"Hours",
		"ID",
	}

	var retainVol []interface{}
	for _, volData := range volArray {
		// create an voluntary profile
		volAttain := model.Vol{
			UserID: profile.ID,
		}
		if volData.ID == "" {
			if result := h.DB.Create(&volAttain); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new vol field.",
				))
			}
		} else {
			if result := h.DB.First(&volAttain, volData.ID); result.Error != nil {
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Volcation Not Found",
				))
			}
		}
		for _, fieldName := range fieldlist {
			fieldValue := reflect.ValueOf(volData).FieldByName(fieldName)
			if fieldValue.IsValid() {
				volfield := model.VolField{
					Key:   fieldName,
					Value: fmt.Sprintf("%v", fieldValue.Interface()),
					VolID: volAttain.ID,
				}
				// Look for VolField
				var oldVolField model.VolField
				if err := h.DB.Where("vol_id = ? AND key = ?", volfield.VolID, volfield.Key).First(&oldVolField).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// If key does not exist, create a new VolField
						if err := h.DB.Create(&volfield).Error; err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
								fiber.StatusInternalServerError,
								"",
								"Error occurred while creating the new vol field.",
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
					// If key exists, update the existing VolField
					if err := h.DB.Model(&oldVolField).Updates(&volfield).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
							fiber.StatusInternalServerError,
							"",
							"Error occurred while updating the vol field.",
						))
					}
				}
			}
		}

		retainVol = append(retainVol, volAttain.ID)
	}

	if err := h.DB.Preload("Voluntaries").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {

		for _, vol := range profile.Voluntaries {
			if !functions.ArrayContains(retainVol, vol.ID) {
				if err := h.DB.Delete(&vol).Error; err != nil {
					fmt.Println("Error Deleting Voluntary", err)
				}
			}
		}

	}

	return nil
}
