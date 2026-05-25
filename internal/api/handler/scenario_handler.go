package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mr-isik/gatling-backend/internal/api/middleware"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/service"
)

type ScenarioHandler struct {
	scenarioService *service.ScenarioService
}

func NewScenarioHandler(scenarioService *service.ScenarioService) *ScenarioHandler {
	return &ScenarioHandler{scenarioService: scenarioService}
}

// List godoc
// @Summary      List Scenarios
// @Description  Lists all test scenarios for a given project.
// @Tags         Scenario
// @Security     BearerAuth
// @Produce      json
// @Param        project_id query string true "Project ID"
// @Success      200  {array}   domain.Scenario
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios [get]
func (h *ScenarioHandler) List(c *fiber.Ctx) error {
	projectID := c.Query("project_id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": domain.ErrBadRequest.Error()})
	}

	scenarios, err := h.scenarioService.List(c.UserContext(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(scenarios)
}

// Create godoc
// @Summary      Create Scenario
// @Description  Creates a new test scenario.
// @Tags         Scenario
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        scenario body domain.Scenario true "Scenario Data"
// @Success      201  {object}  domain.Scenario
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios [post]
func (h *ScenarioHandler) Create(c *fiber.Ctx) error {
	var scenario domain.Scenario
	if err := c.BodyParser(&scenario); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	scenario.CreatedBy = middleware.GetUserIDFromContext(c.UserContext())

	created, err := h.scenarioService.Create(c.UserContext(), &scenario)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// Get godoc
// @Summary      Get Scenario
// @Description  Gets detailed information for a specific scenario by ID.
// @Tags         Scenario
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Scenario ID"
// @Success      200  {object}  domain.Scenario
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios/{id} [get]
func (h *ScenarioHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	scenario, err := h.scenarioService.GetByID(c.UserContext(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(scenario)
}

// Update godoc
// @Summary      Update Scenario
// @Description  Updates the data of an existing scenario by ID.
// @Tags         Scenario
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Scenario ID"
// @Param        scenario body domain.Scenario true "Updated Scenario Data"
// @Success      200  {object}  domain.Scenario
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios/{id} [put]
func (h *ScenarioHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var scenario domain.Scenario
	if err := c.BodyParser(&scenario); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	scenario.ID = id

	updated, err := h.scenarioService.Update(c.UserContext(), &scenario)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

// Delete godoc
// @Summary      Delete Scenario
// @Description  Deletes a scenario from the system by ID.
// @Tags         Scenario
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Scenario ID"
// @Success      204  "No Content"
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios/{id} [delete]
func (h *ScenarioHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scenarioService.Delete(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type generateRequest struct {
	Prompt     string             `json:"prompt"`
	TargetURL  string             `json:"target_url,omitempty"`
	ApiContext *domain.ApiContext `json:"api_context,omitempty"`
}

// Generate godoc
// @Summary      Generate Scenario via AI
// @Description  Generates a test scenario using AI based on a prompt.
// @Tags         Scenario
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body generateRequest true "AI Prompt"
// @Success      200  {object}  domain.Scenario
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios/generate [post]
func (h *ScenarioHandler) Generate(c *fiber.Ctx) error {
	var req generateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	scenario, err := h.scenarioService.Generate(c.UserContext(), req.Prompt, req.TargetURL, req.ApiContext)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(scenario)
}

type cloneRequest struct {
	Name string `json:"name"`
}

// Clone godoc
// @Summary      Clone Scenario
// @Description  Clones an existing scenario with a new name.
// @Tags         Scenario
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Scenario ID to clone"
// @Param        request body cloneRequest true "New Scenario Name"
// @Success      201  {object}  domain.Scenario
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios/{id}/clone [post]
func (h *ScenarioHandler) Clone(c *fiber.Ctx) error {
	id := c.Params("id")
	var req cloneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cloned, err := h.scenarioService.Clone(c.UserContext(), id, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(cloned)
}

// ParseApiDoc godoc
// @Summary      Parse API Documentation
// @Description  Parses uploaded API documentation (OpenAPI/Swagger, cURL, URL) and returns normalized endpoints.
// @Tags         Scenario
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        source formData string true "Source type: openapi, curl, url"
// @Param        file   formData file   false "API doc file (for openapi)"
// @Param        content formData string false "Raw content (for curl/url)"
// @Success      200  {object}  domain.ApiContext
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/scenarios/parse-api-doc [post]
func (h *ScenarioHandler) ParseApiDoc(c *fiber.Ctx) error {
	source := c.FormValue("source")
	parser := service.NewApiDocParser()

	switch source {
	case "openapi":
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file is required for openapi source"})
		}
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		defer file.Close()

		bytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(bytes); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		apiCtx, err := parser.ParseOpenAPISpec(c.UserContext(), bytes)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(apiCtx)

	case "curl":
		content := c.FormValue("content")
		if content == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "content is required for curl source"})
		}
		apiCtx, err := parser.ParseCurlCommands(content)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(apiCtx)

	case "url":
		content := c.FormValue("content")
		if content == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "url content is required for url source"})
		}
		apiCtx, err := parser.FetchAndParseURL(c.UserContext(), content)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(apiCtx)

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid source type"})
	}
}
