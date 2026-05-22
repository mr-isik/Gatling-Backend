package api

import (
	"net/http"

	"github.com/mr-isik/gatling-backend/internal/api/handler"
	"github.com/mr-isik/gatling-backend/internal/api/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/mr-isik/gatling-backend/docs"
)

func SetupRoutes(
	authHandler *handler.AuthHandler,
	scenarioHandler *handler.ScenarioHandler,
	runHandler *handler.TestRunHandler,
	reportHandler *handler.ReportHandler,
	projectHandler *handler.ProjectHandler,
	wsHandler *handler.WSHandler,
	jwtSecret string,
) http.Handler {
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.RefreshToken)

	// Protected Routes setup
	protected := http.NewServeMux()

	// Auth API Keys
	protected.HandleFunc("POST /auth/api-keys", authHandler.CreateAPIKey)
	protected.HandleFunc("DELETE /auth/api-keys/{id}", authHandler.DeleteAPIKey)

	// Scenarios
	protected.HandleFunc("GET /v1/scenarios", scenarioHandler.List)
	protected.HandleFunc("POST /v1/scenarios", scenarioHandler.Create)
	protected.HandleFunc("GET /v1/scenarios/{id}", scenarioHandler.Get)
	protected.HandleFunc("PUT /v1/scenarios/{id}", scenarioHandler.Update)
	protected.HandleFunc("DELETE /v1/scenarios/{id}", scenarioHandler.Delete)
	protected.HandleFunc("POST /v1/scenarios/generate", scenarioHandler.Generate)
	protected.HandleFunc("POST /v1/scenarios/{id}/clone", scenarioHandler.Clone)

	// Test Runs
	protected.HandleFunc("POST /v1/runs", runHandler.Start)
	protected.HandleFunc("GET /v1/runs", runHandler.List)
	protected.HandleFunc("GET /v1/runs/{id}", runHandler.Get)
	protected.HandleFunc("POST /v1/runs/{id}/stop", runHandler.Stop)
	protected.HandleFunc("GET /v1/runs/{id}/metrics", runHandler.Metrics)
	protected.HandleFunc("GET /v1/runs/{id}/logs", runHandler.Logs)

	// Reports
	protected.HandleFunc("GET /v1/reports/{id}", reportHandler.Get)
	protected.HandleFunc("GET /v1/reports/{id}/ai-summary", reportHandler.AISummary)
	protected.HandleFunc("POST /v1/reports/{id}/export", reportHandler.Export)
	protected.HandleFunc("GET /v1/reports/compare", reportHandler.Compare)

	// Projects
	protected.HandleFunc("GET /v1/projects", projectHandler.List)
	protected.HandleFunc("POST /v1/projects", projectHandler.Create)
	protected.HandleFunc("POST /v1/projects/{id}/members", projectHandler.AddMember)
	protected.HandleFunc("GET /v1/projects/{id}/usage", projectHandler.Usage)

	// Wrap protected routes with AuthMiddleware
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	mux.Handle("/", authMiddleware(protected))

	// WS Routes (Public or protected via token query param. For now, public stub)
	mux.HandleFunc("GET /ws/runs/{id}/live", wsHandler.ServeLiveMetrics)
	mux.HandleFunc("GET /ws/runs/{id}/anomalies", wsHandler.ServeAnomalies)

	// Swagger UI
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Global Middleware Chain
	// RateLimit (e.g., 100 req/sec) -> Logger -> CORS -> Mux
	handler := middleware.CORSMiddleware(mux)
	handler = middleware.LoggerMiddleware(handler)
	handler = middleware.RateLimitMiddleware(100)(handler)

	return handler
}
