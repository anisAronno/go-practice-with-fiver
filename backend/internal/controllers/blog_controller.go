package controllers

import (
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/response"
	"gofiver/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BlogController struct {
	blogService *services.BlogService
}

func NewBlogController(blogService *services.BlogService) *BlogController {
	return &BlogController{blogService: blogService}
}

func (ctrl *BlogController) Index(c *fiber.Ctx) error {
	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	blogs, total, err := ctrl.blogService.GetAll(pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch blogs")
	}

	var data []map[string]interface{}
	for _, blog := range blogs {
		data = append(data, blog.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}

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
	currentUser := c.Locals("user").(*models.User)

	blog, err := ctrl.blogService.Update(uint(id), userID, &req, currentUser.IsAdmin())
	if err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(c, "You can only update your own blogs")
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, blog.ToResponse(), "Blog updated successfully")
}

func (ctrl *BlogController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid blog ID")
	}

	userID := c.Locals("userID").(uint)
	currentUser := c.Locals("user").(*models.User)

	if err := ctrl.blogService.Delete(uint(id), userID, currentUser.IsAdmin()); err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(c, "You can only delete your own blogs")
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Blog deleted successfully")
}

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

func (ctrl *BlogController) Trashed(c *fiber.Ctx) error {
	pagination := &dto.PaginationQuery{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 15),
	}

	blogs, total, err := ctrl.blogService.GetAllDeleted(pagination)
	if err != nil {
		return response.InternalError(c, "Failed to fetch deleted blogs")
	}

	var data []map[string]interface{}
	for _, blog := range blogs {
		data = append(data, blog.ToResponse())
	}

	return response.Paginated(c, data, total, pagination.GetPage(), pagination.GetPerPage())
}

func (ctrl *BlogController) Restore(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid blog ID")
	}

	if err := ctrl.blogService.Restore(uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Blog restored successfully")
}

func (ctrl *BlogController) ForceDelete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid blog ID")
	}

	if err := ctrl.blogService.ForceDelete(uint(id)); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, nil, "Blog permanently deleted")
}

func (ctrl *BlogController) UploadImage(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "Invalid blog ID")
	}

	userID := c.Locals("userID").(uint)
	currentUser := c.Locals("user").(*models.User)

	blog, err := ctrl.blogService.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, "Blog not found")
	}

	if blog.UserID != userID && !currentUser.IsAdmin() {
		return response.Forbidden(c, "You can only update your own blogs")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, "Image file required")
	}

	ext := ""
	switch file.Header.Get("Content-Type") {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	default:
		return response.BadRequest(c, "Invalid image format. Use jpg, png, gif, or webp")
	}

	filename := "blog_" + strconv.FormatUint(id, 10) + "_" + strconv.FormatInt(blog.UpdatedAt.Unix(), 10) + ext
	uploadPath := "/uploads/" + filename

	if err := c.SaveFile(file, "/app"+uploadPath); err != nil {
		return response.InternalError(c, "Failed to save image")
	}

	if err := ctrl.blogService.UpdateImage(uint(id), uploadPath); err != nil {
		return response.InternalError(c, "Failed to update blog image")
	}

	return response.Success(c, map[string]string{"image": uploadPath}, "Image uploaded successfully")
}
