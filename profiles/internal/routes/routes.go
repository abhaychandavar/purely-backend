package routes

import (
	"profiles/internal/controllers"
	"profiles/internal/middlewares/authMiddlewares"
	"profiles/internal/providers/storage"
	"profiles/internal/services"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
}

func (r *Router) InitRoutes(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	locationRoutes := LocationRoutes{
		locationController: controllers.LocationController{
			LocationService: services.LocationService{},
		},
	}
	profileRoutes := ProfileRoutes{
		profileController: controllers.ProfileController{
			ProfileService: services.ProfileService{
				StorageProvider: &storage.GCPStorageProvider{},
			},
		},
	}

	// location routes
	locationRoutesGroup := router.Group("/locations")
	locationRoutesGroup.Use(authMiddlewares.VerifyUserAccess)
	locationRoutes.InitRoutes(locationRoutesGroup)

	// profile routes
	profileRoutesGroup := router.Group("/")
	profileRoutesGroup.Use(authMiddlewares.VerifyUserAccess)
	profileRoutes.InitRoutes(profileRoutesGroup)
}
