package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jcalabing/hrmis-api/internal/common/config"
	"github.com/jcalabing/hrmis-api/internal/common/db"
	"github.com/jcalabing/hrmis-api/internal/handler/profile"
	"github.com/jcalabing/hrmis-api/internal/handler/setup"
	"github.com/jcalabing/hrmis-api/internal/handler/token"
	"github.com/jcalabing/hrmis-api/internal/handler/user"
)

func main() {

	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	app := fiber.New()
	app.Use(cors.New())
	db := db.Init(c.DBUrl)

	user.RegisterRoutes(app, db)
	token.RegisterRoutes(app, db)
	profile.RegisterRoutes(app, db)
	setup.RegisterRoutes(app, db)

	app.Listen(c.Port)
}
