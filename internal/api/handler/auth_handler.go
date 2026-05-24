package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-isik/gatling-backend/internal/api/middleware"
	"github.com/mr-isik/gatling-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary      Register User
// @Description  Creates a new user and returns a JWT token pair.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body registerRequest true "User Registration Info"
// @Success      201  {object}  service.TokenPair
// @Failure      400  {object}  map[string]interface{}
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tokenPair, err := h.authService.Register(c.UserContext(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(tokenPair)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary      Login
// @Description  Logs in a user with email and password, returns a JWT token pair.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "Login Credentials"
// @Success      200  {object}  service.TokenPair
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tokenPair, err := h.authService.Login(c.UserContext(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tokenPair)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken godoc
// @Summary      Refresh Token
// @Description  Generates a new access token using a refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body refreshRequest true "Refresh Token"
// @Success      200  {object}  service.TokenPair
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tokenPair, err := h.authService.RefreshToken(c.UserContext(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tokenPair)
}

type createAPIKeyRequest struct {
	Name string `json:"name"`
}

// CreateAPIKey godoc
// @Summary      Create API Key
// @Description  Generates a new API Key for the user (shown only once).
// @Tags         Auth
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body createAPIKeyRequest true "API Key Name"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /auth/api-keys [post]
func (h *AuthHandler) CreateAPIKey(c *fiber.Ctx) error {
	var req createAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := middleware.GetUserIDFromContext(c.UserContext())
	rawKey, keyModel, err := h.authService.CreateAPIKey(c.UserContext(), userID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"raw_key": rawKey,
		"api_key": keyModel,
	})
}

// DeleteAPIKey godoc
// @Summary      Delete API Key
// @Description  Deletes an existing API Key.
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "API Key ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /auth/api-keys/{id} [delete]
func (h *AuthHandler) DeleteAPIKey(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
	}

	if err := h.authService.DeleteAPIKey(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
