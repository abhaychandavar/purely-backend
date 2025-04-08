package authMiddlewares

import (
	"context"
	"log"
	"profiles/internal/config"
	"profiles/internal/types/appTypes"
	firebaseHelper "profiles/internal/utils/helpers/firebaseHelpers"
	httpErrors "profiles/internal/utils/helpers/httpError"
	"profiles/internal/utils/helpers/httpHelper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func VerifyInternalAccess(c *fiber.Ctx) error {
	// Retrieve the 'Access-Token' header
	accessToken := c.Get("Authorization")
	if len(accessToken) == 0 {
		accessToken = c.Get("Access-Token")
	}
	if len(accessToken) == 0 {
		accessToken = c.Query("accessToken")
	}

	if len(accessToken) == 0 {
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/unauthorized", 401, "Unauthorized"))
	}

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
