package server

import (
	"media/internal/config"
	"media/internal/controllers"
	"media/internal/routes"
	"media/internal/services"
	"media/providers/storage"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterFiberRoutes() {
	awsStorageProvider, err := storage.NewAWSStorageProvider(config.GetConfig().AWS.Region, config.GetConfig().AWS.AWSAccessKeyId, config.GetConfig().AWS.AWSSecretAccessKey)
	if err != nil {
		panic(err)
	}
	mediaService := services.MediaService{
		StorageProvider: awsStorageProvider,
	}
	internalRoutesGroup := s.App.Group("/internal")
	internalRoutes := routes.InternalRoutes{
		InternalController: controllers.InternalController{
			MediaService: mediaService,
		},
	}
	internalRoutes.InitRoutes(internalRoutesGroup)

	rootGroup := s.App.Group("/")

	mediaRouter := routes.Router{
		MediaController: controllers.MediaController{
			MediaService: mediaService,
		},
	}
	mediaRouter.InitRoutes(rootGroup)

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
