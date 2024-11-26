package authMiddlewares

import (
	"auth/internal/config"
	"auth/internal/types/appTypes"
	firebaseHelper "auth/internal/utils/helpers/firebaseHelpers"
	httpErrors "auth/internal/utils/helpers/httpError"
	"auth/internal/utils/helpers/httpHelper"
	"context"
	"log"
	"strings"

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

func VerifyUserAccess(c *fiber.Ctx) error {
	authorizationToken := c.Get("Authorization")
	bearerToken := strings.Split(authorizationToken, "Bearer ")
	if len(bearerToken) < 2 {
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/unauthorized", 401, "Unauthorized"))
	}
	firebaseAuth, err := firebaseHelper.App().Auth(context.Background())
	if err != nil {
		log.Default().Println(err)
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/internal_server_error", 500, "Internal Server Error"))
	}
	token, err := firebaseAuth.VerifyIDToken(context.Background(), bearerToken[1])
	if err != nil {
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/unauthorized", 401, "Unauthorized"))
	}
	log.Default().Printf("token %v", token)
	c.Locals("auth", appTypes.Auth{
		Id: token.Claims["id"].(string),
	})
	log.Default().Printf("auth %v", token.Claims)
	return c.Next()
}
