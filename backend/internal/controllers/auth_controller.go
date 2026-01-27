package controllers

import (
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/response"
	"gofiver/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}


	if req.Name == "" || req.Email == "" || req.Password == "" {
		return response.BadRequest(c, "Name, email and password are required")
	}

	if len(req.Password) < 6 {
		return response.BadRequest(c, "Password must be at least 6 characters")
	}

	user, token, err := ctrl.authService.Register(&req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Created(c, fiber.Map{
		"user":  user.ToResponse(),
		"token": token,
	}, "Registration successful")
}

func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return response.BadRequest(c, "Email and password are required")
	}

	user, token, err := ctrl.authService.Login(&req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	return response.Success(c, fiber.Map{
		"user":  user.ToResponse(),
		"token": token,
	}, "Login successful")
}

func (ctrl *AuthController) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return response.Success(c, user.ToResponse())
}
