package profile

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/model"
)

type child struct {
	Fullname    string      `json:"fullname"`
	Dateofbirth string      `json:"dateofbirth"`
	ID          interface{} `json:"id"`
}

func UpdateChildren(h *handler, c *fiber.Ctx, profile model.User, children string) error {
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

	var childArray []child

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(children), &childArray); err != nil {
		tx.Rollback()
		fmt.Println("Error unmarshaling Children:", err)
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"Error occurred while parsing Children JSON.",
		))
	}

	var childIDS []interface{}
	for _, child := range childArray {
		newChild := model.Children{
			Fullname:    child.Fullname,
			Dateofbirth: child.Dateofbirth,
			UserID:      profile.ID,
		}
		if child.ID == "" {
			// Create a new child if ID is empty
			if result := tx.Create(&newChild); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new Children.",
				))
			}
			childIDS = append(childIDS, newChild.ID)
		} else {
			// Update an existing child if ID is not empty
			var oldChildData model.Children
			if result := tx.First(&oldChildData, child.ID); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Education Not Found",
				))
			} else {
				if err := tx.Model(&oldChildData).Updates(&newChild).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating Child Data.",
					))
				}
				childIDS = append(childIDS, oldChildData.ID)
			}
		}
	}

	tx.Commit() // Commit the transaction if all operations are successful
	fmt.Println(childIDS)
	return nil
}
