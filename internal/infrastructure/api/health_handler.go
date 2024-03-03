package api

import (
	"github.com/gofiber/fiber/v2"
)

type HealthResponse struct {
	Version string `json:"version"`
}

func HealthHandler(version string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(HealthResponse{Version: version})
	}
}
