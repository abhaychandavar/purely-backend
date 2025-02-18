package httpHelper

import (
	"context"
	httpErrors "profiles/internal/utils/helpers/httpError"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type successResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
}

type errorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ControllerHelperType struct {
	C             *fiber.Ctx
	DataExtractor func(c *fiber.Ctx) interface{}
	Handler       func(c *context.Context, data interface{}) (interface{}, error)
	Message       *string
	Code          *string
	Ctx           *context.Context
}

func Controller(params ControllerHelperType) error {
	if params.Ctx == nil {
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		params.Ctx = &newCtx
	}
	var extractedData interface{}
	if params.DataExtractor != nil {
		extractedData = params.DataExtractor(params.C)
	}

	data, err := params.Handler(params.Ctx, extractedData)
	if err != nil {
		log.Error(err)
		if httpErr, ok := err.(*httpErrors.HttpError); ok {
			return params.C.Status(httpErr.StatusCode).JSON(errorResponse{
				Code:    httpErr.Code,
				Message: httpErr.Message,
			})
		}
		return params.C.Status(fiber.StatusInternalServerError).JSON(errorResponse{
			Code:    "purely/requests/errors/server",
			Message: "An internal server error occurred",
		})
	}

	// Prepare the message and code for the response
	message := "Request handled successfully"
	if params.Message != nil {
		message = *params.Message
	}

	code := "purely/requests/success"
	if params.Code != nil {
		code = "purely/requests/" + *params.Code
	}

	return params.C.JSON(successResponse{
		Data:    data,
		Message: message,
		Code:    code,
	})
}

func SendErrorResponse(c *fiber.Ctx, err error) error {
	if httpErr, ok := err.(*httpErrors.HttpError); ok {
		return c.Status(httpErr.StatusCode).JSON(errorResponse{
			Code:    httpErr.Code,
			Message: httpErr.Message,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(errorResponse{
		Code:    "purely/requests/errors/server",
		Message: "An internal server error occurred",
	})
}
