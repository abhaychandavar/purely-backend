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
			AllowOrigins:     "http://localhost:3000", // Exact frontend origin
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
			AllowHeaders:     "Content-Type, Authorization, Accept, Accept-Language, Origin, Referer, User-Agent, Sec-CH-UA, Sec-CH-UA-Mobile, Sec-CH-UA-Platform, Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site",
			AllowCredentials: true, // Required if credentials like cookies or Authorization headers are sent
		}))

	}
	if config.GetConfig().Env == "prod" {
		server.App.Use(cors.New(cors.Config{
			AllowOrigins:     "https://purelyapp.me", // Exact frontend origin
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
			AllowHeaders:     "Content-Type, Authorization, Accept, Accept-Language, Origin, Referer, User-Agent, Sec-CH-UA, Sec-CH-UA-Mobile, Sec-CH-UA-Platform, Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site",
			AllowCredentials: true, // Required if credentials like cookies or Authorization headers are sent
		}))

	}

	return server
}
