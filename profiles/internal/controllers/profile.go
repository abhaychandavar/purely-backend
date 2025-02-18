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
)

type ProfileController struct {
	ProfileService services.ProfileService
}

func (this *ProfileController) CreateProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			return this.ProfileService.CreateProfile(ctx, data.(profileServiceTypes.CreateProfileType))
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

func (this *ProfileController) GetProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			selfData, ok := data.(profileServiceTypes.GetProfileType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return this.ProfileService.GetProfile(ctx, selfData)
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

func (this *ProfileController) GetProfileLayout(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			getProfileLayoutData, ok := data.(profileServiceTypes.GetProfileLayoutType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return this.ProfileService.GetProfileLayout(ctx, getProfileLayoutData)
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

func (this *ProfileController) UpsertDatingProfile(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			profileData, ok := data.(profileServiceTypes.UpsertDatingProfileType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}

			// Call the service with the parsed profile data
			response, err := this.ProfileService.UpsertDatingProfile(ctx, &profileData)
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
			log.Printf("Parsed dating profile: %+v", datingProfile)

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
					convertedPrompts = append(convertedPrompts, profileServiceTypes.DatingPromptType{
						PromptId: *prompt.PromptId,
						Answer:   *prompt.Answer,
					})
				}
			}

			// Convert images
			var parsedImages []profileServiceTypes.ImageElementType
			if datingProfile.Images != nil {
				for _, image := range *datingProfile.Images {
					parsedImages = append(parsedImages, profileServiceTypes.ImageElementType{
						ImageId: *image.ImageId,
						Order:   *image.Order,
					})
				}
			}

			// Process location
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
				Images:                 &parsedImages,
				Location:               location,
				LocationLabel:          datingProfile.LocationLabel,
				PreferredMatchDistance: datingProfile.PreferredMatchDistance,
			}
		},
		Message: nil,
		Code:    nil,
	})
}

func (this *ProfileController) GetPrompts(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			getPromptsData, ok := data.(profileServiceTypes.GetPromptsType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return this.ProfileService.GetPrompts(ctx, getPromptsData)
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

func (this *ProfileController) GetGenders(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			// Directly use the extracted data, which is already the expected type.
			getGendersData, ok := data.(profileServiceTypes.GetGendersType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}

			// Call the service with the prepared data.
			return this.ProfileService.GetGenders(ctx, getGendersData)
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

func (this *ProfileController) GetProfiles(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			getProfilesData, ok := data.(profileServiceTypes.GetProfilesType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return this.ProfileService.GetProfiles(ctx, getProfilesData)
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

func (this *ProfileController) GenerateMediaUploadSignedUrl(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			getSignedUrlData, ok := data.(profileServiceTypes.GenerateMediaUploadSignedUrlType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return this.ProfileService.GenerateMediaUploadSignedUrl(ctx, getSignedUrlData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var mediaUploadData profileControllerTypes.GenerateMediaUploadSignedUrlType
			if err := c.BodyParser(&mediaUploadData); err != nil {
				return nil
			}
			auth := c.Locals("auth").(appTypes.Auth)
			authId := auth.Id
			mediaUploadParams := profileServiceTypes.GenerateMediaUploadSignedUrlType{
				FileName: *mediaUploadData.FileName,
				MimeType: *mediaUploadData.MimeType,
				AuthId:   authId,
				FileSize: *mediaUploadData.FileSize,
				Purpose:  *mediaUploadData.Purpose,
			}
			return mediaUploadParams
		},
		Message: nil,
		Code:    nil,
	})
}

func (this *ProfileController) CompleteMediaUpload(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx *context.Context, data interface{}) (interface{}, error) {
			getSignedUrlData, ok := data.(profileServiceTypes.GenerateMediaUploadSignedUrlType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return this.ProfileService.GenerateMediaUploadSignedUrl(ctx, getSignedUrlData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var mediaUploadData profileControllerTypes.GenerateMediaUploadSignedUrlType
			if err := c.BodyParser(&mediaUploadData); err != nil {
				return nil
			}
			auth := c.Locals("auth").(appTypes.Auth)
			authId := auth.Id
			mediaUploadParams := profileServiceTypes.GenerateMediaUploadSignedUrlType{
				FileName: *mediaUploadData.FileName,
				MimeType: *mediaUploadData.MimeType,
				AuthId:   authId,
				FileSize: *mediaUploadData.FileSize,
				Purpose:  *mediaUploadData.Purpose,
			}
			return mediaUploadParams
		},
		Message: nil,
		Code:    nil,
	})
}
