package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/harrisonhys/learn-fiber/config"
	"github.com/harrisonhys/learn-fiber/entities/migration"
	"github.com/harrisonhys/learn-fiber/routers"
)

func main() {
	config.ConnectDB()
	migration.MigrateUser()
	app := fiber.New()
	routers.MainRouter(app)

	app.Listen(":8008")
}
