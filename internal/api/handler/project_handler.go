package handler

import (
	"github.com/gofiber/fiber/v2"
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
func (h *ProjectHandler) List(c *fiber.Ctx) error {
	userID := middleware.GetUserIDFromContext(c.UserContext())
	projects, err := h.projectRepo.List(c.UserContext(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(projects)
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
func (h *ProjectHandler) Create(c *fiber.Ctx) error {
	var project domain.Project
	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	project.OwnerID = middleware.GetUserIDFromContext(c.UserContext())

	created, err := h.projectRepo.Create(c.UserContext(), &project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
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
func (h *ProjectHandler) AddMember(c *fiber.Ctx) error {
	id := c.Params("id")
	var req addMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.projectRepo.AddMember(c.UserContext(), id, req.UserID, req.Role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
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
func (h *ProjectHandler) Usage(c *fiber.Ctx) error {
	// Stub
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_runs": 10,
		"total_vus":  1500,
	})
}
