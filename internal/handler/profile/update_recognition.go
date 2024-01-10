package profile

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
)

type distinction struct {
	Recognition string      `json:"recognition"`
	ID          interface{} `json:"id"`
}

func UpdateRecognitions(h *handler, c *fiber.Ctx, profile model.User, recognition string) error {
	tx := h.DB.Begin()

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"",
			"Error starting database transaction.",
		))
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var recognitionArray []distinction

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(recognition), &recognitionArray); err != nil {
		tx.Rollback()
		fmt.Println("Error unmarshaling Recognition:", err)
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"Error occurred while parsing Recognition JSON.",
		))
	}

	var recognitionIDS []interface{}
	for _, recognition := range recognitionArray {
		newRecognition := model.Recognition{
			Recognition: recognition.Recognition,
			UserID:      profile.ID,
		}
		if recognition.ID == "" {
			// Create a new recognition if ID is empty
			if result := tx.Create(&newRecognition); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new Recognition.",
				))
			}
			recognitionIDS = append(recognitionIDS, newRecognition.ID)
		} else {
			// Update an existing recognition if ID is not empty
			var oldRecognitionData model.Recognition
			if result := tx.First(&oldRecognitionData, recognition.ID); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Education Not Found",
				))
			} else {
				if err := tx.Model(&oldRecognitionData).Updates(&newRecognition).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating Child Data.",
					))
				}
				recognitionIDS = append(recognitionIDS, oldRecognitionData.ID)
			}
		}
	}

	if err := h.DB.Preload("Recognitions").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {
		for _, recognition := range profile.Recognitions {
			if !functions.ArrayContains(recognitionIDS, recognition.ID) {
				if err := h.DB.Delete(&recognition).Error; err != nil {
					tx.Rollback()
					fmt.Println("Error Deleting Education", err)
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while deleting recognition.",
					))
				}
			}
		}

	}

	tx.Commit() // Commit the transaction if all operations are successful
	return nil
}
