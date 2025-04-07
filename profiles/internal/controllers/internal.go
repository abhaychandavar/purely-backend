package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	PubSub "profiles/internal/providers/pubSub"
	"profiles/internal/services"
	"profiles/internal/utils/helpers/httpHelper"

	"github.com/gofiber/fiber/v2"
)

type InternalController struct {
	InternalService services.InternalService
}

type PubSubMessage struct {
	Message struct {
		Data       string            `json:"data"`
		Attributes map[string]string `json:"attributes"`
	} `json:"message"`
}

func (ic *InternalController) HandlePubSubMessage(c *fiber.Ctx) error {
	return httpHelper.Controller(httpHelper.ControllerHelperType{
		C: c,
		Handler: func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Data: ", data)
			res := ic.InternalService.HandlePubSubMessage(ctx, data.(PubSub.PublishMessageType))
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
