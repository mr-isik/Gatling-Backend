package handler

import (
	"net/http"

	"github.com/mr-isik/gatling-backend/internal/api/httputil"
	"github.com/mr-isik/gatling-backend/internal/api/middleware"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type ProjectHandler struct {
	projectRepo repository.ProjectRepository
}

func NewProjectHandler(projectRepo repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepo: projectRepo}
}

// List godoc
// @Summary      List Projects
// @Description  Lists all projects owned by the authenticated user.
// @Tags         Project
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   domain.Project
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/projects [get]
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	projects, err := h.projectRepo.List(r.Context(), userID)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, projects)
}

// Create godoc
// @Summary      Create Project
// @Description  Creates a new project.
// @Tags         Project
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        project body domain.Project true "Project Data"
// @Success      201  {object}  domain.Project
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/projects [post]
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var project domain.Project
	if err := httputil.ReadJSON(r, &project); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	project.OwnerID = middleware.GetUserIDFromContext(r.Context())

	created, err := h.projectRepo.Create(r.Context(), &project)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, created)
}

type addMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// AddMember godoc
// @Summary      Add Project Member
// @Description  Adds a new member to a project.
// @Tags         Project
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Project ID"
// @Param        request body addMemberRequest true "Member Info"
// @Success      200  "OK"
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/projects/{id}/members [post]
func (h *ProjectHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req addMemberRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.projectRepo.AddMember(r.Context(), id, req.UserID, req.Role); err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Usage godoc
// @Summary      Get Project Usage
// @Description  Gets the usage statistics (runs, VUs) for a project.
// @Tags         Project
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Project ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /v1/projects/{id}/usage [get]
func (h *ProjectHandler) Usage(w http.ResponseWriter, r *http.Request) {
	// Stub
	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"total_runs": 10,
		"total_vus":  1500,
	})
}
