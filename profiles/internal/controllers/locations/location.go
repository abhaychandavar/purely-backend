package locationController

import (
	"auth/internal/services/locationService"
	"auth/internal/types/locationControllerTypes"
	"auth/internal/types/locationServiceTypes"
	"auth/internal/utils/helpers/httpHelper"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetLocations(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			log.Default().Printf("data %v", data)
			return locationService.GetLocations(ctx, data.(locationServiceTypes.GetLocationsType))
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
