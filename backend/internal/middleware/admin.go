package middleware

import (
	"gofiver/internal/models"
	"gofiver/internal/response"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		if !user.IsAdmin() {
			return response.Forbidden(c, "Admin access required")
		}
		return c.Next()
	}
}
