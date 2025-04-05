package controllers

import (
	"context"
	"media/internal/services"
	"media/internal/types/appTypes"
	"media/internal/types/mediaControllerTypes"
	"media/internal/types/mediaServiceTypes"
	httpErrors "media/internal/utils/helpers/httpError"
	httpHelper "media/internal/utils/helpers/httpHelper"

	"github.com/gofiber/fiber/v2"
)

type MediaController struct {
	MediaService services.MediaService
}

type blurImageType struct {
	ImageID   string  `json:"imageID"` // It's a good practice to use capital letters for struct fields to make them exportable
	ProfileID *string `json:"profileID"`
}

func (mediaController *MediaController) GenerateMediaUploadSignedUrl(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			getSignedUrlData, ok := data.(mediaServiceTypes.GenerateMediaUploadSignedUrlType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return mediaController.MediaService.GenerateMediaUploadSignedUrl(ctx, getSignedUrlData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var mediaUploadData mediaControllerTypes.GenerateMediaUploadSignedUrlType
			if err := c.BodyParser(&mediaUploadData); err != nil {
				return nil
			}
			auth := c.Locals("auth").(appTypes.Auth)
			authId := auth.Id
			mediaUploadParams := mediaServiceTypes.GenerateMediaUploadSignedUrlType{
				FileName:    *mediaUploadData.FileName,
				ContentType: *mediaUploadData.ContentType,
				AuthId:      authId,
				FileSize:    *mediaUploadData.FileSize,
				Purpose:     *mediaUploadData.Purpose,
			}
			return mediaUploadParams
		},
		Message: nil,
		Code:    nil,
	})
}

func (mediaController *MediaController) GenerateMultipartUploadUrls(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			getSignedUrlData, ok := data.(mediaServiceTypes.GenerateMultipartUploadUrlsType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return mediaController.MediaService.GenerateMultipartUploadUrls(getSignedUrlData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var mediaUploadData mediaControllerTypes.GenerateMultipartMediaUploadSignedUrls
			if err := c.BodyParser(&mediaUploadData); err != nil {
				return nil
			}
			auth := c.Locals("auth").(appTypes.Auth)
			authId := auth.Id
			mediaUploadParams := mediaServiceTypes.GenerateMultipartUploadUrlsType{
				FileName:    mediaUploadData.FileName,
				ContentType: mediaUploadData.ContentType,
				AuthId:      authId,
				FileSize:    mediaUploadData.FileSize,
				Purpose:     mediaUploadData.Purpose,
			}
			return mediaUploadParams
		},
		Message: nil,
		Code:    nil,
	})
}

func (mediaController *MediaController) CompleteMultipartUpload(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			multiPartUploadData, ok := data.(mediaServiceTypes.CompleteMultipartUploadType)
			if !ok {
				return nil, httpErrors.HydrateHttpError("purely/profiles/requests/errors/invalid-data", 400, "Invalid data")
			}
			return mediaController.MediaService.CompleteMultipartUpload(ctx, multiPartUploadData)
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var mediaUploadData mediaControllerTypes.CompleteMultipartUpload
			if err := c.BodyParser(&mediaUploadData); err != nil {
				return nil
			}
			mediaUploadParams := mediaServiceTypes.CompleteMultipartUploadType{
				UploadID: mediaUploadData.UploadID,
				URL:      mediaUploadData.URL,
				Parts:    mediaUploadData.Parts,
			}
			return mediaUploadParams
		},
		Message: nil,
		Code:    nil,
	})
}
