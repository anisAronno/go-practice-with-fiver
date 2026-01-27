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

	api := app.Group("/api")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "message": "GoFiver API", "version": "1.0.0"})
	})

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	api.Get("/blogs", blogController.Index)
	api.Get("/users/:id/blogs", blogController.UserBlogs)

	protected := api.Group("", middleware.AuthMiddleware(cfg, userRepo))
	protected.Get("/auth/me", authController.Me)
	protected.Get("/users", userController.Index)
	protected.Get("/users/:id", userController.Show)
	protected.Put("/users/:id", userController.Update)
	protected.Delete("/users/:id", userController.Delete)
	protected.Get("/blogs/my", blogController.MyBlogs)
	protected.Post("/blogs", blogController.Store)
	protected.Put("/blogs/:id", blogController.Update)
	protected.Delete("/blogs/:id", blogController.Delete)

	api.Get("/blogs/:id", blogController.Show)
}
