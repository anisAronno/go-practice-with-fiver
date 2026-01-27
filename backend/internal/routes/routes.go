
package routes

import (
	"gofiver/internal/config"
	"gofiver/internal/controllers"
	"gofiver/internal/database"
	"gofiver/internal/middleware"
	"gofiver/internal/repositories"
	"gofiver/internal/services"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *database.Database, redis *database.RedisClient, cfg *config.Config) {
	userRepo := repositories.NewUserRepository(db.DB)
	blogRepo := repositories.NewBlogRepository(db.DB)

	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	blogService := services.NewBlogService(blogRepo)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	blogController := controllers.NewBlogController(blogService)

	app.Static("/uploads", "/app/uploads")

	api := app.Group("/api")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "GoFiver API",
			"version": "1.0.0",
		})
	})

	// Public routes (no middleware)
	api.Get("/blogs", blogController.Index)
	api.Get("/blogs/:id", blogController.Show)
	api.Get("/users/:id/blogs", blogController.UserBlogs)

	// Guest-only auth routes (login/register) - blocks authenticated users
	api.Post("/auth/register", middleware.GuestMiddleware(cfg), authController.Register)
	api.Post("/auth/login", middleware.GuestMiddleware(cfg), authController.Login)

	// Protected routes - requires authentication
	auth := middleware.AuthMiddleware(cfg, userRepo)

	api.Get("/auth/me", auth, authController.Me)

	api.Get("/users", auth, userController.Index)
	api.Get("/users/:id", auth, userController.Show)
	api.Put("/users/:id", auth, userController.Update)
	api.Patch("/users/:id/role", auth, userController.UpdateRole)
	api.Delete("/users/:id", auth, userController.Delete)

	api.Get("/blogs/my", auth, blogController.MyBlogs)
	api.Post("/blogs", auth, blogController.Store)
	api.Put("/blogs/:id", auth, blogController.Update)
	api.Delete("/blogs/:id", auth, blogController.Delete)
	api.Post("/blogs/:id/image", auth, blogController.UploadImage)
	api.Delete("/blogs/:id/image", auth, blogController.DeleteImage)

	// Admin routes
	adminAuth := []fiber.Handler{auth, middleware.AdminMiddleware()}

	api.Get("/admin/blogs/trashed", append(adminAuth, blogController.Trashed)...)
	api.Post("/admin/blogs/:id/restore", append(adminAuth, blogController.Restore)...)
	api.Delete("/admin/blogs/:id/force", append(adminAuth, blogController.ForceDelete)...)

	api.Get("/admin/users/trashed", append(adminAuth, userController.Trashed)...)
	api.Post("/admin/users/:id/restore", append(adminAuth, userController.Restore)...)
	api.Delete("/admin/users/:id/force", append(adminAuth, userController.ForceDelete)...)
}