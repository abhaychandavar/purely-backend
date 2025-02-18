package controllers

import (
	"auth/internal/database/models"
	"auth/internal/services"
	httpHelper "auth/internal/utils/helpers/httpHelper"
	"context"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService services.AuthService
}

func (authController *AuthController) InsertAuth(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			return authController.AuthService.InsertAuth(ctx, data.(models.Auth))
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

func (authController *AuthController) GetAuthToken(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			return authController.AuthService.GetAuthToken(ctx, data.(*string))
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			uid := c.Locals("uid").(string)
			return &uid
		},
		Message: nil,
		Code:    nil,
	})
}
