package controllers

import (
	"context"
	"log"
	"profiles/internal/services"
	"profiles/internal/types/appTypes"
	"profiles/internal/types/profileControllerTypes"
	"profiles/internal/types/profileServiceTypes"
	httpErrors "profiles/internal/utils/helpers/httpError"
	httpHelper "profiles/internal/utils/helpers/httpHelper"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileController struct {
	ProfileService services.ProfileService
}

func (provider *ProfileController) CreateProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			return provider.ProfileService.CreateProfile(ctx, data.(profileServiceTypes.CreateProfileType))
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

func (provider *ProfileController) GetProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			selfData, ok := data.(profileServiceTypes.GetProfileType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return provider.ProfileService.GetProfile(ctx, selfData)
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

func (provider *ProfileController) GetProfileLayout(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			getProfileLayoutData, ok := data.(profileServiceTypes.GetProfileLayoutType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return provider.ProfileService.GetProfileLayout(ctx, getProfileLayoutData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			category := c.Params("profileCategory")
			return profileServiceTypes.GetProfileLayoutType{
				Category: &category,
			}
		},
		Message: nil,
		Code:    nil,
	})
}

func (provider *ProfileController) UpsertDatingProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			profileData, ok := data.(profileServiceTypes.UpsertDatingProfileType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}

			// Call the service with the parsed profile data
			response, err := provider.ProfileService.UpsertDatingProfile(ctx, &profileData)
			if err != nil {
				return nil, err
			}
			return response, nil
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var datingProfile profileControllerTypes.UpsertDatingProfileType
			// Parse the request body
			if err := c.BodyParser(&datingProfile); err != nil {
				log.Printf("Error parsing body: %v", err)
				return nil
			}

			if datingProfile.Media != nil {
				log.Printf("Parsed dating profile: %+v", *datingProfile.Media)
			}

			// Extract authenticated user's ID
			auth, ok := c.Locals("auth").(appTypes.Auth)
			if !ok {
				log.Printf("Error extracting auth from context")
				return nil
			}
			authId := auth.Id

			// Convert prompts
			var convertedPrompts []profileServiceTypes.DatingPromptType
			if datingProfile.Prompts != nil {
				for _, prompt := range *datingProfile.Prompts {
					if prompt.PromptId == nil || prompt.Answer == nil {
						continue
					}
					convertedPrompts = append(convertedPrompts, profileServiceTypes.DatingPromptType{
						PromptId: *prompt.PromptId,
						Answer:   *prompt.Answer,
					})
				}
			}

			// Convert images
			var mediaList []profileServiceTypes.MediaElementType
			if datingProfile.Media != nil {
				for _, media := range *datingProfile.Media {
					if media.MediaID == nil || media.Order == nil {
						continue
					}
					blurredMediaRefID := primitive.NewObjectID().Hex()
					mediaList = append(mediaList, profileServiceTypes.MediaElementType{
						MediaID:        *media.MediaID,
						Order:          *media.Order,
						BlurredImageID: blurredMediaRefID,
					})
				}
			}
			var location *profileServiceTypes.Location
			if datingProfile.Location != nil && datingProfile.Location.Lat != nil {
				location = &profileServiceTypes.Location{
					Lat: *datingProfile.Location.Lat,
					Lng: *datingProfile.Location.Lng,
				}
			}
			// Build the service type
			return profileServiceTypes.UpsertDatingProfileType{
				AuthId:                 &authId,
				Name:                   datingProfile.Name,
				Age:                    datingProfile.Age,
				Gender:                 datingProfile.Gender,
				HereFor:                datingProfile.HereFor,
				LookingFor:             datingProfile.LookingFor,
				Bio:                    datingProfile.Bio,
				Prompts:                &convertedPrompts,
				Media:                  &mediaList,
				Location:               location,
				LocationLabel:          datingProfile.LocationLabel,
				PreferredMatchDistance: datingProfile.PreferredMatchDistance,
			}
		},
		Message: nil,
		Code:    nil,
	})
}

func (provider *ProfileController) GetPrompts(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			getPromptsData, ok := data.(profileServiceTypes.GetPromptsType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return provider.ProfileService.GetPrompts(ctx, getPromptsData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			category := c.Params("profileCategory")
			pageStr := c.Query("page", "0")
			page, err := strconv.ParseInt(pageStr, 10, 64)
			if err != nil {
				// Handle invalid "page" query parameter gracefully.
				page = 0
			}
			return profileServiceTypes.GetPromptsType{
				Category: &category,
				Page:     &page,
			}
		},
		Message: nil,
		Code:    nil,
	})
}

func (provider *ProfileController) GetGenders(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			// Directly use the extracted data, which is already the expected type.
			getGendersData, ok := data.(profileServiceTypes.GetGendersType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}

			// Call the service with the prepared data.
			return provider.ProfileService.GetGenders(ctx, getGendersData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			// Extract "page" query parameter and handle parsing errors.
			pageStr := c.Query("page", "0") // Default to "0" if not provided.
			page, err := strconv.ParseInt(pageStr, 10, 64)
			if err != nil {
				// Handle invalid "page" query parameter gracefully.
				page = 0
			}

			// Return the expected type directly.
			return profileServiceTypes.GetGendersType{
				Page: &page,
			}
		},
		Message: nil,
		Code:    nil,
	})
}

func (provider *ProfileController) GetProfiles(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			getProfilesData, ok := data.(profileServiceTypes.GetProfilesType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return provider.ProfileService.GetProfiles(ctx, getProfilesData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			pageStr := c.Query("page", "0")
			page, err := strconv.ParseInt(pageStr, 10, 64)

			category := c.Params("profileCategory")

			if err != nil {
				page = 0
			}

			auth := c.Locals("auth").(appTypes.Auth)
			authId := auth.Id

			return profileServiceTypes.GetProfilesType{
				Page:     &page,
				Category: category,
				AuthId:   authId,
			}
		},
		Message: nil,
		Code:    nil,
	})
}
