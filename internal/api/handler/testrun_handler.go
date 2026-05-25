package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mr-isik/gatling-backend/internal/api/middleware"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/service"
)

type TestRunHandler struct {
	runService *service.TestRunService
}

func NewTestRunHandler(runService *service.TestRunService) *TestRunHandler {
	return &TestRunHandler{runService: runService}
}

type startRunRequest struct {
	ScenarioID string           `json:"scenario_id"`
	ProjectID  string           `json:"project_id"`
	Config     domain.RunConfig `json:"config"`
}

// Start godoc
// @Summary      Start Test Run
// @Description  Starts a new load test execution.
// @Tags         TestRun
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body startRunRequest true "Start Run Details"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/runs [post]
func (h *TestRunHandler) Start(c *fiber.Ctx) error {
	var req startRunRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := middleware.GetUserIDFromContext(c.UserContext())
	runID, err := h.runService.Start(c.UserContext(), req.ScenarioID, req.ProjectID, userID, req.Config)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": runID})
}

// List godoc
// @Summary      List Test Runs
// @Description  Lists all test runs for a given project.
// @Tags         TestRun
// @Security     BearerAuth
// @Produce      json
// @Param        project_id query string true "Project ID"
// @Success      200  {array}   domain.TestRun
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/runs [get]
func (h *TestRunHandler) List(c *fiber.Ctx) error {
	projectID := c.Query("project_id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": domain.ErrBadRequest.Error()})
	}

	runs, err := h.runService.List(c.UserContext(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(runs)
}

// Get godoc
// @Summary      Get Test Run
// @Description  Gets detailed information for a specific test run by ID.
// @Tags         TestRun
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Test Run ID"
// @Success      200  {object}  domain.TestRun
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/runs/{id} [get]
func (h *TestRunHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	run, err := h.runService.GetByID(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(run)
}

// Stop godoc
// @Summary      Stop Test Run
// @Description  Stops a currently running test run.
// @Tags         TestRun
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Test Run ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/runs/{id}/stop [post]
func (h *TestRunHandler) Stop(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.runService.Stop(c.UserContext(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "stopping"})
}

// Metrics godoc
// @Summary      Get Run Metrics
// @Description  Gets the time-series metrics for a test run.
// @Tags         TestRun
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Test Run ID"
// @Param        from query string false "Start Time (RFC3339)"
// @Param        to query string false "End Time (RFC3339)"
// @Success      200  {array}   domain.AggregatedMetric
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/runs/{id}/metrics [get]
func (h *TestRunHandler) Metrics(c *fiber.Ctx) error {
	id := c.Params("id")

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		from = time.Now().Add(-1 * time.Hour) // default last 1 hour
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		to = time.Now()
	}

	metrics, err := h.runService.GetMetrics(c.UserContext(), id, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(metrics)
}

func (h *TestRunHandler) Logs(c *fiber.Ctx) error {
	// Stub for logs
	return c.Status(fiber.StatusOK).JSON([]string{"log1", "log2"})
}
