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

func RequireAuth(DB *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(viper.GetString("SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
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
}
