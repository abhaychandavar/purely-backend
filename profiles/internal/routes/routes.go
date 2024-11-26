package routes

import (
	"auth/internal/middlewares/authMiddlewares"
	locationRoutes "auth/internal/routes/location"
	profileRoutes "auth/internal/routes/profile"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// location routes
	locationRoutesGroup := router.Group("/locations")
	locationRoutesGroup.Use(authMiddlewares.VerifyUserAccess)
	locationRoutes.InitRoutes(locationRoutesGroup)

	// profile routes
	profileRoutesGroup := router.Group("/")
	profileRoutesGroup.Use(authMiddlewares.VerifyUserAccess)
	profileRoutes.InitRoutes(profileRoutesGroup)
}
