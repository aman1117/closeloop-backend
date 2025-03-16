package routes

import (
	"closeloop/controllers"
	middleware "closeloop/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/users", controllers.RegisterUser)
	app.Post("/login", controllers.LoginUser)
	app.Get("/users", middleware.AuthMiddleware, controllers.GetUser)
}
