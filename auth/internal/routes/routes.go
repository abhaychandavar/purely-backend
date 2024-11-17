package routes

import (
	"auth/internal/controllers/authController"
	"auth/internal/middlewares/authMiddlewares"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	router.Post("/internal", authMiddlewares.VerifyInternalAccess, authController.InsertAuth)
	router.Get("/token", authMiddlewares.ValidateFirebaseToken, authController.GetAuthToken)
}
