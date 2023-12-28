package profile

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
)

type assoc struct {
	Organization string      `json:"organization"`
	ID           interface{} `json:"id"`
}

func UpdateAssociations(h *handler, c *fiber.Ctx, profile model.User, associations string) error {
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

	var assocArray []assoc

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(associations), &assocArray); err != nil {
		tx.Rollback()
		fmt.Println("Error unmarshaling Associations:", err)
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"Error occurred while parsing Associations JSON.",
		))
	}

	var assocIDS []interface{}
	for _, assoc := range assocArray {
		newAssoc := model.Assoc{
			Organization: assoc.Organization,
			UserID:       profile.ID,
		}
		if assoc.ID == "" {
			// Create a new assoc if ID is empty
			if result := tx.Create(&newAssoc); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new Associations.",
				))
			}
			assocIDS = append(assocIDS, newAssoc.ID)
		} else {
			// Update an existing assoc if ID is not empty
			var oldAssocData model.Assoc
			if result := tx.First(&oldAssocData, assoc.ID); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Education Not Found",
				))
			} else {
				if err := tx.Model(&oldAssocData).Updates(&newAssoc).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating Assoc Data.",
					))
				}
				assocIDS = append(assocIDS, oldAssocData.ID)
			}
		}
	}

	if err := h.DB.Preload("Associations").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {
		for _, assoc := range profile.Associations {
			if !functions.ArrayContains(assocIDS, assoc.ID) {
				if err := h.DB.Delete(&assoc).Error; err != nil {
					tx.Rollback()
					fmt.Println("Error Deleting Education", err)
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while deleting assoc.",
					))
				}
			}
		}

	}

	tx.Commit() // Commit the transaction if all operations are successful
	return nil
}
