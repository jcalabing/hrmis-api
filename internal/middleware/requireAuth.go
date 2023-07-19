package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jcalabing/hrmis-api/internal/common/errors"
	"github.com/jcalabing/hrmis-api/internal/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func RequireAuth(DB *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		// tokenString := c.Cookies("Authorization")
		// return c.Status(fiber.StatusOK).SendString(tokenString)
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
			if !ok || float64(time.Now().Unix()) > claims["exp"].(float64) {
				return c.Status(fiber.StatusUnauthorized).SendString("Token has expired")
			}

			var profile model.User
			if err := DB.First(&profile, claims["sub"]).Error; err != nil {
				return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
			}

			c.Locals("profile", profile)

			return c.Next()
		}

		//if all fails unauthorized
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
}
