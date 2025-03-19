package controllers

import (
	"context"
	"log"
	"profiles/internal/services"
	"profiles/internal/types/locationControllerTypes"
	"profiles/internal/types/locationServiceTypes"
	"profiles/internal/utils/helpers/httpHelper"

	"github.com/gofiber/fiber/v2"
)

type LocationController struct {
	LocationService services.LocationService
}

func (locationController *LocationController) GetLocations(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			log.Default().Printf("data %v", data)
			return locationController.LocationService.GetLocations(ctx, data.(locationServiceTypes.GetLocationsType))
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var locationBody locationControllerTypes.GetLocationsType
			if err := c.QueryParser(&locationBody); err != nil {
				return nil
			}
			return locationServiceTypes.GetLocationsType{
				Query:     &locationBody.Query,
				PageToken: &locationBody.PageToken,
			}
		},
		Message: nil,
		Code:    nil,
	})
}
