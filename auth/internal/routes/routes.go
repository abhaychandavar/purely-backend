package routes

import (
	"auth/internal/controllers"
	"auth/internal/middlewares/authMiddlewares"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	AuthController controllers.AuthController
}

func (r *Router) InitRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	router.Post("/internal", authMiddlewares.VerifyInternalAccess, r.AuthController.InsertAuth)
	router.Get("/token", authMiddlewares.ValidateFirebaseToken, r.AuthController.GetAuthToken)
}
