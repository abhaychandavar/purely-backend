package routes

import (
	"profiles/internal/config"
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

	storageProviderInstance, err := storage.NewAWSStorageProvider(config.GetConfig().AWS.Region, config.GetConfig().AWS.AWSAccessKeyId, config.GetConfig().AWS.AWSSecretAccessKey)
	if err != nil {
		panic(err)
	}

	profileRoutes := ProfileRoutes{
		profileController: controllers.ProfileController{
			ProfileService: services.ProfileService{
				StorageProvider: storageProviderInstance,
			},
		},
	}

	locationRoutesGroup := router.Group("/locations")
	locationRoutesGroup.Use(authMiddlewares.VerifyUserAccess)
	locationRoutes.InitRoutes(locationRoutesGroup)

	profileRoutesGroup := router.Group("/")
	profileRoutesGroup.Use(authMiddlewares.VerifyUserAccess)
	profileRoutes.InitRoutes(profileRoutesGroup)
}
