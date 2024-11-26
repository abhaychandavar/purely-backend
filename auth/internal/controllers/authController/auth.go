package authController

import (
	"auth/internal/database/models"
	authService "auth/internal/services"
	httpHelper "auth/internal/utils/helpers/httpHelper"
	"context"

	"github.com/gofiber/fiber/v2"
)

func InsertAuth(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			return authService.InsertAuth(ctx, data.(models.Auth))
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var auth models.Auth
			if err := c.BodyParser(&auth); err != nil {
				return nil
			}

			return auth
		},
		Message: nil,
		Code:    nil,
	})
}

func GetAuthToken(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			return authService.GetAuthToken(ctx, data.(*string))
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			uid := c.Locals("uid").(string)
			return &uid
		},
		Message: nil,
		Code:    nil,
	})
}
