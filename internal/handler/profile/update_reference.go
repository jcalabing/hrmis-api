package profile

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
)

type reference struct {
	Fullname  string      `json:"fullname"`
	Address   string      `json:"address"`
	Telephone string      `json:"telephone"`
	ID        interface{} `json:"id"`
}

func UpdateReferences(h *handler, c *fiber.Ctx, profile model.User, references string) error {
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

	var referenceArray []reference

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(references), &referenceArray); err != nil {
		tx.Rollback()
		fmt.Println("Error unmarshaling References:", err)
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"Error occurred while parsing References JSON.",
		))
	}

	var referenceIDS []interface{}
	for _, reference := range referenceArray {
		newRefence := model.Reference{
			Fullname:  reference.Fullname,
			Address:   reference.Address,
			Telephone: reference.Telephone,
			UserID:    profile.ID,
		}
		if reference.ID == "" {
			// Create a new reference if ID is empty
			if result := tx.Create(&newRefence); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new References.",
				))
			}
			referenceIDS = append(referenceIDS, newRefence.ID)
		} else {
			// Update an existing reference if ID is not empty
			var oldRefenceData model.Reference
			if result := tx.First(&oldRefenceData, reference.ID); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Education Not Found",
				))
			} else {
				if err := tx.Model(&oldRefenceData).Updates(&newRefence).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating Refence Data.",
					))
				}
				referenceIDS = append(referenceIDS, oldRefenceData.ID)
			}
		}
	}

	if err := h.DB.Preload("References").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {
		for _, reference := range profile.References {
			if !functions.ArrayContains(referenceIDS, reference.ID) {
				if err := h.DB.Delete(&reference).Error; err != nil {
					tx.Rollback()
					fmt.Println("Error Deleting Education", err)
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while deleting reference.",
					))
				}
			}
		}

	}

	tx.Commit() // Commit the transaction if all operations are successful
	return nil
}
