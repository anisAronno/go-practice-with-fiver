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
	svc *services.BlogService
}

func NewBlogController(svc *services.BlogService) *BlogController {
	return &BlogController{svc: svc}
}

func (c *BlogController) Index(ctx *fiber.Ctx) error {
	q := &dto.BlogSearchQuery{
		Page:    ctx.QueryInt("page", 1),
		PerPage: ctx.QueryInt("per_page", 20),
		Search:  ctx.Query("search"),
	}

	blogs, total, err := c.svc.GetAll(q)
	if err != nil {
		return response.InternalError(ctx, "Failed to fetch blogs")
	}

	data := make([]map[string]interface{}, len(blogs))
	for i, b := range blogs {
		data[i] = map[string]interface{}{
			"id":         b.ID,
			"title":      b.Title,
			"image":      b.Image,
			"user_id":    b.UserID,
			"created_at": b.CreatedAt,
			"user":       map[string]interface{}{"id": b.UserID, "name": b.UserName},
		}
	}

	return response.Paginated(ctx, data, total, q.GetPage(), q.GetPerPage())
}

func (c *BlogController) Show(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid blog ID")
	}

	blog, err := c.svc.GetByID(uint(id))
	if err != nil {
		return response.NotFound(ctx, "Blog not found")
	}

	return response.Success(ctx, blog.ToResponse())
}

func (c *BlogController) Store(ctx *fiber.Ctx) error {
	var req dto.CreateBlogRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid request body")
	}
	if req.Title == "" || req.Content == "" {
		return response.BadRequest(ctx, "Title and content are required")
	}

	userID := ctx.Locals("userID").(uint)
	blog, err := c.svc.Create(userID, &req)
	if err != nil {
		return response.InternalError(ctx, "Failed to create blog")
	}

	return response.Created(ctx, blog.ToResponse(), "Blog created successfully")
}

func (c *BlogController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid blog ID")
	}

	var req dto.UpdateBlogRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid request body")
	}

	userID := ctx.Locals("userID").(uint)
	user := ctx.Locals("user").(*models.User)

	blog, err := c.svc.Update(uint(id), userID, &req, user.IsAdmin())
	if err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(ctx, "You can only update your own blogs")
		}
		return response.BadRequest(ctx, err.Error())
	}

	return response.Success(ctx, blog.ToResponse(), "Blog updated successfully")
}

func (c *BlogController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid blog ID")
	}

	userID := ctx.Locals("userID").(uint)
	user := ctx.Locals("user").(*models.User)

	if err := c.svc.Delete(uint(id), userID, user.IsAdmin()); err != nil {
		if err.Error() == "unauthorized" {
			return response.Forbidden(ctx, "You can only delete your own blogs")
		}
		return response.BadRequest(ctx, err.Error())
	}

	return response.Success(ctx, nil, "Blog deleted successfully")
}

func (c *BlogController) MyBlogs(ctx *fiber.Ctx) error {
	q := &dto.PaginationQuery{
		Page:    ctx.QueryInt("page", 1),
		PerPage: ctx.QueryInt("per_page", 20),
	}

	userID := ctx.Locals("userID").(uint)
	blogs, total, err := c.svc.GetByUserID(userID, q)
	if err != nil {
		return response.InternalError(ctx, "Failed to fetch blogs")
	}

	data := make([]map[string]interface{}, len(blogs))
	for i, b := range blogs {
		data[i] = b.ToListResponse()
	}

	return response.Paginated(ctx, data, total, q.GetPage(), q.GetPerPage())
}

func (c *BlogController) UserBlogs(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid user ID")
	}

	q := &dto.PaginationQuery{
		Page:    ctx.QueryInt("page", 1),
		PerPage: ctx.QueryInt("per_page", 20),
	}

	blogs, total, err := c.svc.GetByUserID(uint(userID), q)
	if err != nil {
		return response.InternalError(ctx, "Failed to fetch blogs")
	}

	data := make([]map[string]interface{}, len(blogs))
	for i, b := range blogs {
		data[i] = b.ToListResponse()
	}

	return response.Paginated(ctx, data, total, q.GetPage(), q.GetPerPage())
}

func (c *BlogController) Trashed(ctx *fiber.Ctx) error {
	q := &dto.PaginationQuery{
		Page:    ctx.QueryInt("page", 1),
		PerPage: ctx.QueryInt("per_page", 20),
	}

	blogs, total, err := c.svc.GetAllDeleted(q)
	if err != nil {
		return response.InternalError(ctx, "Failed to fetch deleted blogs")
	}

	data := make([]map[string]interface{}, len(blogs))
	for i, b := range blogs {
		data[i] = b.ToListResponse()
	}

	return response.Paginated(ctx, data, total, q.GetPage(), q.GetPerPage())
}

func (c *BlogController) Restore(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid blog ID")
	}

	if err := c.svc.Restore(uint(id)); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	return response.Success(ctx, nil, "Blog restored successfully")
}

func (c *BlogController) ForceDelete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid blog ID")
	}

	if err := c.svc.ForceDelete(uint(id)); err != nil {
		return response.BadRequest(ctx, err.Error())
	}

	return response.Success(ctx, nil, "Blog permanently deleted")
}

func (c *BlogController) UploadImage(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid blog ID")
	}

	userID := ctx.Locals("userID").(uint)
	user := ctx.Locals("user").(*models.User)

	blog, err := c.svc.GetByID(uint(id))
	if err != nil {
		return response.NotFound(ctx, "Blog not found")
	}

	if blog.UserID != userID && !user.IsAdmin() {
		return response.Forbidden(ctx, "You can only update your own blogs")
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return response.BadRequest(ctx, "Image file required")
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
		return response.BadRequest(ctx, "Invalid image format")
	}

	filename := "blog_" + strconv.FormatUint(id, 10) + "_" + strconv.FormatInt(blog.UpdatedAt.Unix(), 10) + ext
	uploadPath := "/uploads/" + filename

	if err := ctx.SaveFile(file, "/app"+uploadPath); err != nil {
		return response.InternalError(ctx, "Failed to save image")
	}

	if err := c.svc.UpdateImage(uint(id), uploadPath); err != nil {
		return response.InternalError(ctx, "Failed to update blog image")
	}

	return response.Success(ctx, map[string]string{"image": uploadPath}, "Image uploaded")
}

func (c *BlogController) DeleteImage(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.BadRequest(ctx, "Invalid blog ID")
	}

	userID := ctx.Locals("userID").(uint)
	user := ctx.Locals("user").(*models.User)

	blog, err := c.svc.GetByID(uint(id))
	if err != nil {
		return response.NotFound(ctx, "Blog not found")
	}

	if blog.UserID != userID && !user.IsAdmin() {
		return response.Forbidden(ctx, "You can only update your own blogs")
	}

	if err := c.svc.UpdateImage(uint(id), ""); err != nil {
		return response.InternalError(ctx, "Failed to delete image")
	}

	return response.Success(ctx, nil, "Image deleted")
}
