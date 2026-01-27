package controllers

import (
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/response"
	"gofiver/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) Index(c *fiber.Ctx) error {
	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	users, total, err := ctrl.userService.GetAll(pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch users")
	}

	var data []map[string]interface{}
	for _, user := range users {
		data = append(data, user.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}

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

func (ctrl *UserController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	currentUser := c.Locals("user").(*models.User)
	if currentUser.ID != uint(id) && !currentUser.IsAdmin() {
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

func (ctrl *UserController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	currentUser := c.Locals("user").(*models.User)
	if currentUser.ID != uint(id) && !currentUser.IsAdmin() {
		return response.Forbidden(c, "You can only delete your own account")
	}

	if currentUser.ID == uint(id) && currentUser.IsAdmin() {
		return response.Forbidden(c, "Admin cannot delete themselves")
	}

	if err := ctrl.userService.Delete(uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "User deleted successfully")
}

func (ctrl *UserController) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	currentUser := c.Locals("user").(*models.User)
	if !currentUser.IsAdmin() {
		return response.Forbidden(c, "Only admin can change roles")
	}

	if currentUser.ID == uint(id) {
		return response.Forbidden(c, "Cannot change your own role")
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Role != "admin" && req.Role != "author" {
		return response.BadRequest(c, "Role must be 'admin' or 'author'")
	}

	user, err := ctrl.userService.UpdateRole(uint(id), req.Role)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, user.ToResponse(), "Role updated successfully")
}

func (ctrl *UserController) Trashed(c *fiber.Ctx) error {
	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	users, total, err := ctrl.userService.GetAllDeleted(pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch deleted users")
	}

	var data []map[string]interface{}
	for _, user := range users {
		data = append(data, user.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}

func (ctrl *UserController) Restore(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := ctrl.userService.Restore(uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "User restored successfully")
}

func (ctrl *UserController) ForceDelete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	if err := ctrl.userService.ForceDelete(uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "User permanently deleted")
}
