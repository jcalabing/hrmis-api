package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jcalabing/hrmis-api/internal/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const UserContextKey = "user"

func RequireAuth(DB *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("Authorization")

		if tokenString != "" {

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte(viper.GetString("SECRET")), nil
			})

			if err != nil {
				return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
				if time.Now().After(expirationTime) {
					return c.Status(fiber.StatusUnauthorized).SendString("Token has expired")
				}

				var profile model.User
				if err := DB.First(&profile, claims["sub"]).Error; err != nil {
					return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
				}

				c.Locals("profile", profile)

				return c.Next()
			}

			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")

		}

		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
}
