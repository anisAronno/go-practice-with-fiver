
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

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "GoFiver API",
			"version": "1.0.0",
		})
	})

	guest := api.Group("/auth", middleware.GuestMiddleware(cfg))
	guest.Post("/register", authController.Register)
	guest.Post("/login", authController.Login)

	api.Get("/blogs", blogController.Index)
	api.Get("/blogs/:id", blogController.Show)
	api.Get("/users/:id/blogs", blogController.UserBlogs)

	protected := api.Group("", middleware.AuthMiddleware(cfg, userRepo))
	
	protected.Get("/auth/me", authController.Me)
	
	protected.Get("/users", userController.Index)
	protected.Get("/users/:id", userController.Show)
	protected.Put("/users/:id", userController.Update)
	protected.Patch("/users/:id/role", userController.UpdateRole)
	protected.Delete("/users/:id", userController.Delete)
	
	protected.Get("/blogs/my", blogController.MyBlogs)
	protected.Post("/blogs", blogController.Store)
	protected.Put("/blogs/:id", blogController.Update)
	protected.Delete("/blogs/:id", blogController.Delete)
	protected.Post("/blogs/:id/image", blogController.UploadImage)

	admin := api.Group("/admin", middleware.AuthMiddleware(cfg, userRepo), middleware.AdminMiddleware())
	
	admin.Get("/blogs/trashed", blogController.Trashed)
	admin.Post("/blogs/:id/restore", blogController.Restore)
	admin.Delete("/blogs/:id/force", blogController.ForceDelete)
	
	admin.Get("/users/trashed", userController.Trashed)
	admin.Post("/users/:id/restore", userController.Restore)
	admin.Delete("/users/:id/force", userController.ForceDelete)
}