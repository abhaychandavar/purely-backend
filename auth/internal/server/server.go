package server

import (
	"github.com/gofiber/fiber/v2"

	"auth/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "auth",
			AppName:      "auth",
		}),

		db: database.Mongo(),
	}
	return server
}
