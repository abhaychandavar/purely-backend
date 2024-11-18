package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"auth/internal/config"
	"auth/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "profiles",
			AppName:      "profiles",
		}),

		db: database.Mongo(),
	}

	if config.GetConfig().Env != "prod" {
		server.App.Use(cors.New(cors.Config{
			AllowOrigins: "http://localhost:3000",       // Only allow requests from this origin
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Specify allowed HTTP methods
			AllowHeaders: "Content-Type, Authorization", // Specify allowed headers
		}))
	}
	if config.GetConfig().Env == "prod" {
		server.App.Use(cors.New(cors.Config{
			AllowOrigins: "https://purelyapp.me",        // Only allow requests from this origin
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Specify allowed HTTP methods
			AllowHeaders: "Content-Type, Authorization", // Specify allowed headers
		}))
	}

	return server
}
