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

type eli struct {
	Eligibility string      `json:"eligibility"`
	Rating      string      `json:"rating"`
	Dated       string      `json:"dated"`
	Exam        string      `json:"exam"`
	Number      string      `json:"number"`
	Validity    string      `json:"validity"`
	ID          interface{} `json:"id"`
}

func UpdateEli(h *handler, c *fiber.Ctx, profile model.User, eligibility string) error {
	var eliArray []eli

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(eligibility), &eliArray); err != nil {
		fmt.Println("Error unmarshaling Eligibility:", err)
	}

	fieldlist := []string{
		"Eligibility",
		"Rating",
		"Dated",
		"Exam",
		"Number",
		"Validity",
		"ID",
	}

	var retainEli []interface{}
	for _, eliData := range eliArray {
		// create an eligibility profile
		eliAttain := model.Eli{
			UserID: profile.ID,
		}
		if eliData.ID == "" {
			if result := h.DB.Create(&eliAttain); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new eli field.",
				))
			}
		} else {
			if result := h.DB.First(&eliAttain, eliData.ID); result.Error != nil {
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Eligibility Not Found",
				))
			}
		}
		for _, fieldName := range fieldlist {
			fieldValue := reflect.ValueOf(eliData).FieldByName(fieldName)
			if fieldValue.IsValid() {
				elifield := model.EliField{
					Key:   fieldName,
					Value: fmt.Sprintf("%v", fieldValue.Interface()),
					EliID: eliAttain.ID,
				}
				// Look for EliField
				var oldEliField model.EliField
				if err := h.DB.Where("eli_id = ? AND key = ?", elifield.EliID, elifield.Key).First(&oldEliField).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// If key does not exist, create a new EliField
						if err := h.DB.Create(&elifield).Error; err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
								fiber.StatusInternalServerError,
								"",
								"Error occurred while creating the new eli field.",
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
					// If key exists, update the existing EliField
					if err := h.DB.Model(&oldEliField).Updates(&elifield).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
							fiber.StatusInternalServerError,
							"",
							"Error occurred while updating the eli field.",
						))
					}
				}
			}
		}

		retainEli = append(retainEli, eliAttain.ID)
	}

	if err := h.DB.Preload("Eligibilities").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {

		for _, eli := range profile.Eligibilities {
			if !functions.ArrayContains(retainEli, eli.ID) {
				if err := h.DB.Delete(&eli).Error; err != nil {
					fmt.Println("Error Deleting Eligibility", err)
				}
			}
		}

	}

	return nil
}
