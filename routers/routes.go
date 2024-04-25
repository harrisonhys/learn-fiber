package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/harrisonhys/learn-fiber/controllers"
)

func MainRouter(c *fiber.App) {
	c.Get("/api/users", controllers.UserControllerShow)
	c.Get("/api/users/:userId", controllers.UserControllerById)
	c.Post("/api/users", controllers.UserControllerCreate)
	c.Put("/api/users/:userId", controllers.UserControllerUpdate)
}
