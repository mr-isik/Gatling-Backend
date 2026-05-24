package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mr-isik/gatling-backend/internal/api/httputil"
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
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	tokenPair, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tokenPair)
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
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	tokenPair, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httputil.JSONError(w, http.StatusUnauthorized, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenPair)
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
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	tokenPair, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		httputil.JSONError(w, http.StatusUnauthorized, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokenPair)
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
func (h *AuthHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	var req createAPIKeyRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	rawKey, keyModel, err := h.authService.CreateAPIKey(r.Context(), userID, req.Name)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, map[string]interface{}{
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
func (h *AuthHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httputil.JSONError(w, http.StatusBadRequest, nil)
		return
	}

	if err := h.authService.DeleteAPIKey(r.Context(), id); err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
