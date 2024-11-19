package authMiddlewares

import (
	"auth/internal/config"
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

func ValidateFirebaseToken(c *fiber.Ctx) error {
	authorizationToken := c.Get("Authorization")
	if authorizationToken == "" {
		log.Printf("[ERROR] Missing Authorization header in request: %s %s", c.Method(), c.Path())
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/unauthorized", 401, "Missing Authorization header"))
	}

	bearerToken := strings.Split(authorizationToken, "Bearer ")
	if len(bearerToken) != 2 {
		log.Printf("[ERROR] Invalid Authorization header format: %s %s. Token is missing or malformed", c.Method(), c.Path())
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/unauthorized", 401, "Invalid Authorization header format"))
	}

	// Log Firebase app config or relevant context for debugging
	app := firebaseHelper.App()
	if app == nil {
		log.Printf("[ERROR] Firebase app is nil. Make sure the Firebase app is initialized properly.")
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/internal_server_error", 500, "Internal Server Error"))
	}

	// Try to initialize Firebase Auth and log the failure
	firebaseAuth, err := app.Auth(context.Background())
	if err != nil {
		log.Printf("[ERROR] Firebase Auth initialization failed. Error: %v", err) // Log the error and config
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/internal_server_error", 500, "Internal Server Error"))
	}

	// Verify the Firebase ID Token
	token, err := firebaseAuth.VerifyIDToken(context.Background(), bearerToken[1])
	if err != nil {
		log.Printf("[ERROR] Failed to verify Firebase token: %v, Request: %s %s", err, c.Method(), c.Path())
		return httpHelper.SendErrorResponse(c, httpErrors.HydrateHttpError("purely/requests/errors/unauthorized", 401, "Unauthorized"))
	}

	// Store UID in locals for subsequent use
	c.Locals("uid", token.UID)

	// Proceed to next middleware or handler
	return c.Next()
}
