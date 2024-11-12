package authMiddlewares

import (
	"auth/internal/config"
	httpErrors "auth/internal/utils/helpers/httpError"
	"auth/internal/utils/helpers/httpHelper"

	"github.com/gofiber/fiber/v2"
)

func VerifyInternalAccess(c *fiber.Ctx) error {
	// Retrieve the 'Access-Token' header
	accessToken := c.Get("Access-Token")

	// Check if the token matches the expected internal access token
	if config.GetConfig().InternalAccessToken != accessToken {
		// Send an unauthorized error response if the token does not match
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/unauthorized", 401, "Unauthorized"))
	}

	// Proceed to the next middleware/handler
	return c.Next()
}
