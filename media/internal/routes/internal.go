package routes

import (
	"media/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

type InternalRoutes struct {
	InternalController controllers.InternalController
}

func (ir *InternalRoutes) InitRoutes(router fiber.Router) {
	router.Post("/images/blur", ir.InternalController.BlurImage)
	router.Post("/pubsub/messages", ir.InternalController.HandlePubSubMessage)
}
