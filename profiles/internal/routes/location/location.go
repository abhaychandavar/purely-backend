package locationRoutes

import (
	locationController "auth/internal/controllers/locations"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(router fiber.Router) {
	router.Get("/", locationController.GetLocations)
}
