package profileRoutes

import (
	profileController "auth/internal/controllers/authController"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(router fiber.Router) {
	router.Get("/:profileCategory/self", profileController.GetProfile)
	router.Post("/", profileController.CreateProfile)
}
