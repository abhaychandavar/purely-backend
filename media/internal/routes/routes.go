package routes

import (
	"media/internal/controllers"
	"media/internal/middlewares/authMiddlewares"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	MediaController controllers.MediaController
}

func (r *Router) InitRoutes(router fiber.Router) {
	mediaRouteGroup := router.Group("/")
	mediaRouteGroup.Use(authMiddlewares.VerifyUserAccess)
	mediaRouteGroup.Post("/media/multipart/complete", r.MediaController.CompleteMultipartUpload)
	mediaRouteGroup.Post("/media/multipart/signed-urls", r.MediaController.GenerateMultipartUploadUrls)
}
