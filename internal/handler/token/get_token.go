package token

import (
	// "fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/model"
	"github.com/spf13/viper"
)

func (h *handler) GetToken(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")

	// fmt.Println(authorization)
	// Check if the Authorization header is present
	if authorization == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())
	}

	// Check if the Authorization header starts with "Bearer "
	if !strings.HasPrefix(authorization, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())
	}

	// Extract the token by removing the "Bearer " prefix
	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	// Validate and process the token as needed
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(viper.GetString("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//find the user based on ID
		// fmt.Println(claims["exp"], claims["sub"])
		var user model.User
		h.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())

		}

		returnValue := map[string]any{
			"user": model.ConvertToUserResponse(h.DB, user),
		}
		return c.Status(fiber.StatusOK).JSON(returnValue)
	}

	//return error on invalid Authorization
	return c.Status(fiber.StatusUnauthorized).JSON(errors.UnAuthorized())

}
