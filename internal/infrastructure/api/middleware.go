package api

import "github.com/gofiber/fiber/v2"

type Middleware interface {
	Handler(ctx *fiber.Ctx) error
}
