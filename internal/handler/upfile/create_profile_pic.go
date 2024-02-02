package upfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/model"
	"gorm.io/gorm"
)

func (h *handler) CreateProfilePic(c *fiber.Ctx) error {
	file, err := c.FormFile("picture")
	userID := c.FormValue("userid")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"",
			err,
		))
	}

	// Validate the file extension
	fileExtension := filepath.Ext(file.Filename)
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExtensions[fileExtension] {
		// Return an error response if the file extension is not allowed
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"",
			"Invalid filetype",
		))
	}

	now := time.Now()
	nsec := strconv.FormatInt(now.UnixNano(), 10)

	if userID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
			fiber.StatusInternalServerError,
			"",
			"User Not Found",
		))
	} else {
		var profile model.User
		if result := h.DB.First(&profile, userID); result.Error != nil {
			// return fiber.NewError(fiber.StatusNotFound, result.Error.Error())
			return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
				fiber.StatusBadRequest,
				"",
				"User Not Found",
			))
		}
		userfield := model.UserField{
			Key:    "profilepic",
			Value:  nsec + fileExtension,
			UserID: profile.ID,
		}

		var oldUserField model.UserField
		if err := h.DB.Where("user_id = ? AND key = ?", userfield.UserID, userfield.Key).First(&oldUserField).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// If key does not exist, create a new UserField
				if err := h.DB.Create(&userfield).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
						fiber.StatusInternalServerError,
						"",
						"Error occurred while creating the new user field.",
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
			// If key exists, update the existing UserField
			//delete file
			if err := os.Remove("public/profile/" + oldUserField.Value); err != nil {
				// Handle error if deletion fails
				fmt.Println("Error deleting file:", err)
			}
			if err := h.DB.Model(&oldUserField).Updates(&userfield).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(errors.NewErrorResponse(
					fiber.StatusInternalServerError,
					"",
					"Error occurred while updating the user field.",
				))
			}
		}

		c.SaveFile(file, "public/profile/"+userfield.Value)
	}
	return c.Status(fiber.StatusOK).JSON("OK upfile")
}
