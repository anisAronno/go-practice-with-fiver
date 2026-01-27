package controllers

import (
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/response"
	"gofiver/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// UserController handles user CRUD endpoints
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// Index lists all users
// GET /api/users
func (ctrl *UserController) Index(c *fiber.Ctx) error {
	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	users, total, err := ctrl.userService.GetAll(pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch users")
	}

	// Convert to response format
	var data []map[string]interface{}
	for _, user := range users {
		data = append(data, user.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}

// Show returns a single user
// GET /api/users/:id
func (ctrl *UserController) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	user, err := ctrl.userService.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.Success(c, user.ToResponse())
}

// Update updates a user (self only)
// PUT /api/users/:id
func (ctrl *UserController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	// Check if updating self
	currentUser := c.Locals("user").(*models.User)
	if currentUser.ID != uint(id) {
		return response.Forbidden(c, "You can only update your own profile")
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	user, err := ctrl.userService.Update(uint(id), &req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, user.ToResponse(), "User updated successfully")
}

// Delete deletes a user (self only)
// DELETE /api/users/:id
func (ctrl *UserController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	// Check if deleting self
	currentUser := c.Locals("user").(*models.User)
	if currentUser.ID != uint(id) {
		return response.Forbidden(c, "You can only delete your own account")
	}

	if err := ctrl.userService.Delete(uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "User deleted successfully")
}
