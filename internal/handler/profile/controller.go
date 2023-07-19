package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/middleware"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	routes := app.Group("/profile", middleware.RequireAuth(h.DB))
	routes.Get("/", h.GetProfile)
	routes.Post("/", h.UpdateProfile)
}
