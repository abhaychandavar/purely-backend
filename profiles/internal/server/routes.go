package server

import (
	"profiles/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	locationRoutes := routes.Router{}
	rootGroup := s.App.Group("/")
	locationRoutes.InitRoutes(rootGroup)

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
