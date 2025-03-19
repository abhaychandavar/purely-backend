package routes

import (
	"profiles/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

type ProfileRoutes struct {
	profileController controllers.ProfileController
}

func (profileRoutes *ProfileRoutes) InitRoutes(router fiber.Router) {
	router.Get("/genders", profileRoutes.profileController.GetGenders)
	router.Get("/prompts/:profileCategory", profileRoutes.profileController.GetPrompts)
	router.Get("/:profileCategory", profileRoutes.profileController.GetProfile)
	router.Post("/", profileRoutes.profileController.CreateProfile)
	router.Get("/:profileCategory/layout", profileRoutes.profileController.GetProfileLayout)
	router.Patch("/:profileCategory/upsert", profileRoutes.profileController.UpsertDatingProfile)
	router.Get("/:profileCategory/profiles", profileRoutes.profileController.GetProfiles)
	router.Post("/media/signed-url", profileRoutes.profileController.GenerateMediaUploadSignedUrl)
	router.Post("/media/multipart/signed-urls", profileRoutes.profileController.GenerateMultipartUploadUrls)
	router.Post("/media/multipart/complete", profileRoutes.profileController.CompleteMultipartUpload)
}
