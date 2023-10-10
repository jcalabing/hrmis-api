package token

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/model"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type CreateTokenRequestBody struct {
	Username string `json:"username" validate:"required,nonempty"`
	Password string `json:"password" validate:"required,nonempty"`
}

func (h *handler) CreateToken(c *fiber.Ctx) error {
	body := new(CreateTokenRequestBody)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"",
		))
	}

	// Look for User
	var user model.User
	h.DB.Where("username = ? AND active =?", body.Username, "YES").First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())

	}

	//compare sent in password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {

		return c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(viper.GetString("SECRET")))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.NewErrorResponse(
			fiber.StatusBadRequest,
			"",
			"Failed to create Token.",
		))
	}

	convertedUser := model.ConvertToUserResponse(h.DB, user)

	returnValue := map[string]any{
		"token": tokenString,
		"user":  convertedUser,
	}

	return c.Status(fiber.StatusOK).JSON(returnValue)
}
