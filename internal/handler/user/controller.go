package user

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

	routes := app.Group("/user", middleware.RequireAuth(h.DB))
	routes.Get("/:id", h.GetUser)
	routes.Get("/", h.GetUsers)
	routes.Post("/", h.CreateUser)
	routes.Put("/:id", h.UpdateUser)
	routes.Patch("/:id", h.UpdatePass)

}
