package handler

import (
	"net/http"
	"time"

	"github.com/mr-isik/gatling-backend/internal/api/httputil"
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
func (h *TestRunHandler) Start(w http.ResponseWriter, r *http.Request) {
	var req startRunRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	runID, err := h.runService.Start(r.Context(), req.ScenarioID, req.ProjectID, userID, req.Config)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, map[string]string{"run_id": runID})
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
func (h *TestRunHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		httputil.JSONError(w, http.StatusBadRequest, domain.ErrBadRequest)
		return
	}

	runs, err := h.runService.List(r.Context(), projectID)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, runs)
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
func (h *TestRunHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	run, err := h.runService.GetByID(r.Context(), id)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, run)
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
func (h *TestRunHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.runService.Stop(r.Context(), id); err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "stopping"})
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
func (h *TestRunHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		from = time.Now().Add(-1 * time.Hour) // default last 1 hour
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		to = time.Now()
	}

	metrics, err := h.runService.GetMetrics(r.Context(), id, from, to)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, metrics)
}

func (h *TestRunHandler) Logs(w http.ResponseWriter, r *http.Request) {
	// Stub for logs
	httputil.JSON(w, http.StatusOK, []string{"log1", "log2"})
}
