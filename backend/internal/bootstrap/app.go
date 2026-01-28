package bootstrap

import (
	"fmt"
	"gofiver/internal/config"
	"gofiver/internal/database"
	"gofiver/internal/database/seeders"
	"gofiver/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Application struct {
	Config *config.Config
	Fiber  *fiber.App
	DB     *database.Database
	Redis  *database.RedisClient
}

func NewApplication() (*Application, error) {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName:               cfg.App.Name,
		ErrorHandler:          customErrorHandler,
		DisableStartupMessage: cfg.App.Env == "production",
		Prefork:               false,
		ReduceMemoryUsage:     true, 
	})

	app.Use(recover.New())

	if cfg.App.Debug {
		app.Use(logger.New())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := db.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if err := db.CreateIndexes(); err != nil {
		fmt.Printf("Warning: index creation had issues: %v\n", err)
	}

	redis, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	application := &Application{Config: cfg, Fiber: app, DB: db, Redis: redis}
	routes.RegisterRoutes(app, db, redis, cfg)

	return application, nil
}

func (a *Application) RunSeeders() error {
	return seeders.NewSeeder(a.DB.DB).Run()
}

func (a *Application) Start() error {
	return a.Fiber.Listen(":" + a.Config.App.Port)
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{"success": false, "message": message})
}
