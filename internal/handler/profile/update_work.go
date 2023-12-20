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

type work struct {
	Started    string      `json:"started"`
	Ended      string      `json:"ended"`
	Title      string      `json:"title"`
	Department string      `json:"department"`
	Salary     string      `json:"salary"`
	Grade      string      `json:"grade"`
	Status     string      `json:"status"`
	Gov        string      `json:"gov"`
	ID         interface{} `json:"id"`
}

func UpdateWork(h *handler, c *fiber.Ctx, profile model.User, works string) error {
	var workArray []work

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(works), &workArray); err != nil {
		fmt.Println("Error unmarshaling Works:", err)
	}

	fieldlist := []string{
		"Started",
		"Ended",
		"Title",
		"Department",
		"Salary",
		"Grade",
		"Status",
		"Gov",
		"ID",
	}

	var retainWork []interface{}
	for _, workData := range workArray {
		// create an work profile
		workAttain := model.Work{
			UserID: profile.ID,
		}
		if workData.ID == "" {
			if result := h.DB.Create(&workAttain); result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new work field.",
				))
			}
		} else {
			if result := h.DB.First(&workAttain, workData.ID); result.Error != nil {
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Work Not Found",
				))
			}
		}
		for _, fieldName := range fieldlist {
			fieldValue := reflect.ValueOf(workData).FieldByName(fieldName)
			if fieldValue.IsValid() {
				workfield := model.WorkField{
					Key:    fieldName,
					Value:  fmt.Sprintf("%v", fieldValue.Interface()),
					WorkID: workAttain.ID,
				}
				// Look for WorkField
				var oldWorkField model.WorkField
				if err := h.DB.Where("work_id = ? AND key = ?", workfield.WorkID, workfield.Key).First(&oldWorkField).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// If key does not exist, create a new WorkField
						if err := h.DB.Create(&workfield).Error; err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
								fiber.StatusInternalServerError,
								"",
								"Error occurred while creating the new work field.",
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
					// If key exists, update the existing WorkField
					if err := h.DB.Model(&oldWorkField).Updates(&workfield).Error; err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
							fiber.StatusInternalServerError,
							"",
							"Error occurred while updating the work field.",
						))
					}
				}
			}
		}

		retainWork = append(retainWork, workAttain.ID)
	}

	if err := h.DB.Preload("Works").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {

		for _, work := range profile.Works {
			if !functions.ArrayContains(retainWork, work.ID) {
				if err := h.DB.Delete(&work).Error; err != nil {
					fmt.Println("Error Deleting Work", err)
				}
			}
		}

	}

	return nil
}
