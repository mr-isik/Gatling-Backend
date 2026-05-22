package handler

import (
	"net/http"

	"github.com/mr-isik/gatling-backend/internal/api/httputil"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/service"
)

type ReportHandler struct {
	reportService   *service.ReportService
	baselineService *service.BaselineService
}

func NewReportHandler(reportService *service.ReportService, baselineService *service.BaselineService) *ReportHandler {
	return &ReportHandler{
		reportService:   reportService,
		baselineService: baselineService,
	}
}

// Get godoc
// @Summary      Get Report
// @Description  Gets a specific test report by ID.
// @Tags         Report
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Report ID"
// @Success      200  {object}  domain.Report
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/reports/{id} [get]
func (h *ReportHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	report, err := h.reportService.Get(r.Context(), id)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, report)
}

// AISummary godoc
// @Summary      Get AI Summary
// @Description  Gets an AI generated summary for a specific test report.
// @Tags         Report
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Report ID"
// @Success      200  {object}  domain.Report
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/reports/{id}/ai-summary [get]
func (h *ReportHandler) AISummary(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	report, err := h.reportService.AISummary(r.Context(), id)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, report)
}

type exportRequest struct {
	Format string `json:"format"`
}

// Export godoc
// @Summary      Export Report
// @Description  Exports a report in the specified format (e.g. PDF, CSV).
// @Tags         Report
// @Security     BearerAuth
// @Accept       json
// @Produce      application/octet-stream
// @Param        id path string true "Report ID"
// @Param        request body exportRequest true "Export Format"
// @Success      200  {file}    []byte
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/reports/{id}/export [post]
func (h *ReportHandler) Export(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req exportRequest
	if err := httputil.ReadJSON(r, &req); err != nil {
		httputil.JSONError(w, http.StatusBadRequest, err)
		return
	}

	data, err := h.reportService.Export(r.Context(), id, req.Format)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	// Just a stub for downloading
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Compare godoc
// @Summary      Compare Reports
// @Description  Compares two test runs (baseline vs current).
// @Tags         Report
// @Security     BearerAuth
// @Produce      json
// @Param        run1 query string true "Baseline Run ID"
// @Param        run2 query string true "Current Run ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /v1/reports/compare [get]
func (h *ReportHandler) Compare(w http.ResponseWriter, r *http.Request) {
	run1 := r.URL.Query().Get("run1")
	run2 := r.URL.Query().Get("run2")

	if run1 == "" || run2 == "" {
		httputil.JSONError(w, http.StatusBadRequest, domain.ErrBadRequest)
		return
	}

	comp, err := h.baselineService.CompareWithBaseline(r.Context(), run1, run2)
	if err != nil {
		httputil.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	httputil.JSON(w, http.StatusOK, comp)
}
