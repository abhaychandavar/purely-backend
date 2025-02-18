package server

import (
	"auth/internal/controllers"
	"auth/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	rootGroup := s.App.Group("/")

	authRouter := routes.Router{
		AuthController: controllers.AuthController{},
	}
	authRouter.InitRoutes(rootGroup)

	s.App.Get("/health", s.healthHandler)

}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
