package profile

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/common/functions"
	"github.com/jcalabing/hrmis-api/internal/model"
)

type skill struct {
	Skillhobby  string      `json:"skillhobby"`
	Recognition string      `json:"recognition"`
	ID          interface{} `json:"id"`
}

func UpdateSkills(h *handler, c *fiber.Ctx, profile model.User, skillren string) error {
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

	var skillArray []skill

	// Unmarshal the JSON array into the map
	if err := json.Unmarshal([]byte(skillren), &skillArray); err != nil {
		tx.Rollback()
		fmt.Println("Error unmarshaling Skill:", err)
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"Error occurred while parsing Skill JSON.",
		))
	}

	var skillIDS []interface{}
	for _, skill := range skillArray {
		newSkill := model.Skill{
			Skillhobby:  skill.Skillhobby,
			Recognition: skill.Recognition,
			UserID:      profile.ID,
		}
		if skill.ID == "" {
			// Create a new skill if ID is empty
			if result := tx.Create(&newSkill); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while creating the new Skill.",
				))
			}
			skillIDS = append(skillIDS, newSkill.ID)
		} else {
			// Update an existing skill if ID is not empty
			var oldSkillData model.Skill
			if result := tx.First(&oldSkillData, skill.ID); result.Error != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
					fiber.StatusBadRequest,
					"",
					"Education Not Found",
				))
			} else {
				if err := tx.Model(&oldSkillData).Updates(&newSkill).Error; err != nil {
					tx.Rollback()
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while updating Child Data.",
					))
				}
				skillIDS = append(skillIDS, oldSkillData.ID)
			}
		}
	}

	if err := h.DB.Preload("Skills").Find(&profile).Error; err != nil {
		fmt.Println("Error Retrieving Data: ", err)
	} else {
		for _, skill := range profile.Skills {
			if !functions.ArrayContains(skillIDS, skill.ID) {
				if err := h.DB.Delete(&skill).Error; err != nil {
					tx.Rollback()
					fmt.Println("Error Deleting Education", err)
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while deleting skill.",
					))
				}
			}
		}

	}

	tx.Commit() // Commit the transaction if all operations are successful
	return nil
}
