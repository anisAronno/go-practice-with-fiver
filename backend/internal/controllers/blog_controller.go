package controllers

import (
	"gofiver/internal/dto"
	"gofiver/internal/response"
	"gofiver/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// BlogController handles blog CRUD endpoints
type BlogController struct {
	blogService *services.BlogService
}

// NewBlogController creates a new blog controller
func NewBlogController(blogService *services.BlogService) *BlogController {
	return &BlogController{blogService: blogService}
}

// Index lists all blogs
// GET /api/blogs
func (ctrl *BlogController) Index(c *fiber.Ctx) error {
	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	blogs, total, err := ctrl.blogService.GetAll(pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch blogs")
	}

	// Convert to response format
	var data []map[string]interface{}
	for _, blog := range blogs {
		data = append(data, blog.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}

// Show returns a single blog
// GET /api/blogs/:id
func (ctrl *BlogController) Show(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid blog ID")
	}

	blog, err := ctrl.blogService.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, "Blog not found")
	}

	return response.Success(c, blog.ToResponse())
}

// Store creates a new blog
// POST /api/blogs
func (ctrl *BlogController) Store(c *fiber.Ctx) error {
	var req dto.CreateBlogRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Title == "" || req.Content == "" {
		return response.BadRequest(c, "Title and content are required")
	}

	userID := c.Locals("userID").(uint)

	blog, err := ctrl.blogService.Create(userID, &req)
	if err != nil {
		return response.InternalError(c, "Failed to create blog")
	}

	return response.Created(c, blog.ToResponse(), "Blog created successfully")
}

// Update updates a blog
// PUT /api/blogs/:id
func (ctrl *BlogController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid blog ID")
	}

	var req dto.UpdateBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	userID := c.Locals("userID").(uint)

	blog, err := ctrl.blogService.Update(uint(id), userID, &req)
	if err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(c, "You can only update your own blogs")
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, blog.ToResponse(), "Blog updated successfully")
}

// Delete deletes a blog
// DELETE /api/blogs/:id
func (ctrl *BlogController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid blog ID")
	}

	userID := c.Locals("userID").(uint)

	if err := ctrl.blogService.Delete(uint(id), userID); err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(c, "You can only delete your own blogs")
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Blog deleted successfully")
}

// MyBlogs returns blogs of authenticated user
// GET /api/blogs/my
func (ctrl *BlogController) MyBlogs(c *fiber.Ctx) error {
	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	userID := c.Locals("userID").(uint)

	blogs, total, err := ctrl.blogService.GetByUserID(userID, pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch blogs")
	}

	var data []map[string]interface{}
	for _, blog := range blogs {
		data = append(data, blog.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}

// UserBlogs returns blogs of a specific user
// GET /api/users/:id/blogs
func (ctrl *BlogController) UserBlogs(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID")
	}

	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	blogs, total, err := ctrl.blogService.GetByUserID(uint(userID), pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch blogs")
	}

	var data []map[string]interface{}
	for _, blog := range blogs {
		data = append(data, blog.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}
