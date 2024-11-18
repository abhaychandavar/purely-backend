package routes

import (
	"auth/internal/middlewares/authMiddlewares"
	profileRoutes "auth/internal/routes/profile"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// profile routes
	profileRoutesGroup := router.Group("/profiles")
	profileRoutesGroup.Use(authMiddlewares.VerifyUserAccess)
	profileRoutes.InitRoutes(profileRoutesGroup)
}
