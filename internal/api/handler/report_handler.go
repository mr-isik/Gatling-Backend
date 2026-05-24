package handler

import (
	"github.com/gofiber/fiber/v2"
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
func (h *ReportHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	report, err := h.reportService.Get(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(report)
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
func (h *ReportHandler) AISummary(c *fiber.Ctx) error {
	id := c.Params("id")
	report, err := h.reportService.AISummary(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(report)
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
func (h *ReportHandler) Export(c *fiber.Ctx) error {
	id := c.Params("id")
	var req exportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	data, err := h.reportService.Export(c.UserContext(), id, req.Format)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Just a stub for downloading
	c.Set("Content-Type", "application/octet-stream")
	return c.Status(fiber.StatusOK).Send(data)
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
func (h *ReportHandler) Compare(c *fiber.Ctx) error {
	run1 := c.Query("run1")
	run2 := c.Query("run2")

	if run1 == "" || run2 == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": domain.ErrBadRequest.Error()})
	}

	comp, err := h.baselineService.CompareWithBaseline(c.UserContext(), run1, run2)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(comp)
}
