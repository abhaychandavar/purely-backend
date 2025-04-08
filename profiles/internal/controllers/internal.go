package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	PubSub "profiles/internal/providers/pubSub"
	"profiles/internal/services"
	"profiles/internal/utils/helpers/httpHelper"

	"github.com/gofiber/fiber/v2"
)

type PubSubMessagePayload struct {
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
}

type InternalController struct {
	InternalService services.InternalService
}

func (ic *InternalController) HandlePubSubMessage(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			res := ic.InternalService.HandlePubSubMessage(ctx, data.(PubSub.PubSubMessageType))
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
			var data PubSub.PubSubMessageType
			if err := json.Unmarshal(decoded, &data); err != nil {
				fmt.Println("Failed to parse JSON from decoded message: ", err)
				return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON payload")
			}
			return data
		},
		Message: nil,
		Code:    nil,
	})
}
