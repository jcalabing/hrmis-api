package setup

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jcalabing/hrmis-api/internal/model"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	// routes := app.Group("/setup")
	// routes.Get("/", h.GetSetup)
	// routes.Post("/", h.UpdateSetup)

	usersExist := h.UsersExist()

	// Conditional route registration based on user existence
	if usersExist {
		routes := app.Group("/setup")
		routes.Get("/", h.GetSetup)
		routes.Post("/", h.UpdateSetup)
	} else {
		// Add alternative routes or handle the case where users don't exist
		// For example, you can add a route to create an initial user, etc.
	}
}

func (h *handler) UsersExist() bool {
	var count int64
	var user model.User
	h.DB.Model(&user).Count(&count)
	return count == 0
}
