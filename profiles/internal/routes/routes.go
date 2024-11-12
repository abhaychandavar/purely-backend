package routes

import (
	profileController "auth/internal/controllers/authController"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	router.Post("/", profileController.CreateProfile)
}
