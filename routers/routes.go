package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/harrisonhys/learn-fiber/controllers"
)

func MainRouter(c *fiber.App) {
	c.Get("/api/users", controllers.UserControllerShow)
	c.Post("/api/users", controllers.UserControllerCreate)
}
