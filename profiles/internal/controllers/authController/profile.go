package profileController

import (
	profileService "auth/internal/services"
	"auth/internal/types/appTypes"
	"auth/internal/types/profileControllerTypes"
	"auth/internal/types/profileServiceTypes"
	httpErrors "auth/internal/utils/helpers/httpError"
	httpHelper "auth/internal/utils/helpers/httpHelper"

	"github.com/gofiber/fiber/v2"
)

func CreateProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(data interface{}) (interface{}, error) {
			return profileService.CreateProfile(data.(profileServiceTypes.CreateProfileType))
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var profile profileControllerTypes.CreateProfileType
			if err := c.BodyParser(&profile); err != nil {
				return nil
			}

			auth := c.Locals("auth").(appTypes.Auth)
			authId := auth.Id

			userProfile := profileServiceTypes.CreateProfileType{
				AuthId: &authId,
				Lat:    profile.Lat,
				Lng:    profile.Lng,
			}
			return userProfile
		},
		Message: nil,
		Code:    nil,
	})
}

func GetProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(data interface{}) (interface{}, error) {
			selfData, ok := data.(profileServiceTypes.GetProfileType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return profileService.GetProfile(selfData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			auth := c.Locals("auth").(appTypes.Auth)
			authId := auth.Id

			category := c.Params("profileCategory")

			return profileServiceTypes.GetProfileType{
				AuthId:   &authId,
				Category: &category,
			}
		},
		Message: nil,
		Code:    nil,
	})
}
