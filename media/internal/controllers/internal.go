package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"media/internal/services"
	"media/internal/utils/helpers/httpHelper"
	PubSub "media/providers/pubSub"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type InternalController struct {
	MediaService services.MediaService
}

type PubSubMessagePayload struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}

type PubSubMessageType struct {
	Data map[string]interface{} `json:"data"`
	Type string                 `json:"type"`
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
			fmt.Println("HandlePubSubMessage Data: ", data)
			parsedData, ok := data.(PubSub.PublishMessageType)
			if !ok {
				fmt.Println("HandlePubSubMessage Type parsing error")
				return false, nil
			}
			res := ic.MediaService.HandlePubSubMessage(ctx, parsedData)
			return res, nil
		},
		DataExtractor: func(c *fiber.Ctx) interface{} {
			var msg PubSubMessagePayload
			if err := c.BodyParser(&msg); err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
			}
			decoded, err := base64.StdEncoding.DecodeString(msg.Message.Data)
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Failed to decode message data"})
			}
			var data map[string]interface{}
			if err := json.Unmarshal(decoded, &data); err != nil {
				fmt.Println("Failed to parse JSON from decoded message: ", err)
				return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON payload")
			}
			fmt.Println("HandlePubSubMessage Data: ", data)
			return data
		},
		Message: nil,
		Code:    nil,
	})
}
