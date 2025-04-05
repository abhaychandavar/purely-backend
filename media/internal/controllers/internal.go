package controllers

import (
	"context"
	"encoding/json"
	"media/internal/services"
	"media/internal/utils/helpers/httpHelper"
	PubSub "media/providers/pubSub"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type InternalController struct {
	MediaService services.MediaService
}

type PubSubMessage struct {
	Message struct {
		Data       string            `json:"data"`
		Attributes map[string]string `json:"attributes"`
	} `json:"message"`
}

func (ic *InternalController) BlurImage(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			profileID := data.(blurImageType).ProfileID
			res, err := ic.MediaService.BlurImage(ctx, data.(blurImageType).ImageID, profileID)
			if err != nil {
				return nil, err
			}
			return *res, nil
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var data blurImageType
			if err := c.BodyParser(&data); err != nil {
				return nil
			}
			return data
		},
		Message: nil,
		Code:    nil,
	})
}

func (ic *InternalController) HandlePubSubMessage(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			res := ic.MediaService.HandlePubSubMessage(ctx, data.(PubSub.PublishMessageType))
			return res, nil
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var msg PubSubMessage
			if err := c.BodyParser(&msg); err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
			}
			var pubSubMessageData PubSub.PublishMessageType
			if err := json.Unmarshal([]byte(msg.Message.Data), &pubSubMessageData); err != nil {
				return map[string]interface{}{}
			}
			return pubSubMessageData
		},
		Message: nil,
		Code:    nil,
	})
}
