package handler

import (
	"net/http"

	"github.com/mr-isik/gatling-backend/internal/api/httputil"
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
func (h *ScenarioHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		httputil.JSONError(w, http.StatusBadRequest, domain.ErrBadRequest)
		return
	}

	scenarios, err := h.scenarioService.List(r.Context(), projectID)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, scenarios)
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
func (h *ScenarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var scenario domain.Scenario
	if err := httputil.ReadJSON(r, &scenario); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	scenario.CreatedBy = middleware.GetUserIDFromContext(r.Context())

	created, err := h.scenarioService.Create(r.Context(), &scenario)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, created)
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
func (h *ScenarioHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	scenario, err := h.scenarioService.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			httputil.JSONError(w, http.StatusNotFound, err)
		} else {
			httputil.JSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	httputil.JSON(w, http.StatusOK, scenario)
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
func (h *ScenarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var scenario domain.Scenario
	if err := httputil.ReadJSON(r, &scenario); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}
	scenario.ID = id

	updated, err := h.scenarioService.Update(r.Context(), &scenario)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, updated)
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
func (h *ScenarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.scenarioService.Delete(r.Context(), id); err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type generateRequest struct {
	Prompt string `json:"prompt"`
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
func (h *ScenarioHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	scenario, err := h.scenarioService.Generate(r.Context(), req.Prompt)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, scenario)
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
func (h *ScenarioHandler) Clone(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req cloneRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	cloned, err := h.scenarioService.Clone(r.Context(), id, req.Name)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusCreated, cloned)
}
