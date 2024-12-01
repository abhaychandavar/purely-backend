package profileRoutes

import (
	profileController "auth/internal/controllers/profilesController"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(router fiber.Router) {
	router.Get("/genders", profileController.GetGenders)
	router.Get("/prompts/:profileCategory", profileController.GetPrompts)
	router.Get("/:profileCategory", profileController.GetProfile)
	router.Post("/", profileController.CreateProfile)
	router.Get("/:profileCategory/layout", profileController.GetProfileLayout)
	router.Patch("/:profileCategory/upsert", profileController.UpsertDatingProfile)
	router.Get("/:profileCategory/profiles", profileController.GetProfiles)
}
