package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response provides consistent API responses (like Laravel's API Resources)

// Success returns a success response
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// Created returns a 201 created response
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// Paginated returns a paginated response
func Paginated(c *fiber.Ctx, data interface{}, total int64, page, perPage int) error {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     perPage,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// Error returns an error response
func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}

// BadRequest returns a 400 response
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

// Unauthorized returns a 401 response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

// Forbidden returns a 403 response
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message)
}

// NotFound returns a 404 response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

// InternalError returns a 500 response
func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}
