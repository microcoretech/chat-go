package api

import "github.com/gofiber/fiber/v2"

type Controller interface {
	SetupRoutes(r fiber.Router)
}
