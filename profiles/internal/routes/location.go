package routes

import (
	"profiles/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

type LocationRoutes struct {
	locationController controllers.LocationController
}

func (locationRouter *LocationRoutes) InitRoutes(router fiber.Router) {
	router.Get("/", locationRouter.locationController.GetLocations)
}
