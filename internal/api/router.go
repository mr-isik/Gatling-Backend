package api

import (
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	_ "github.com/mr-isik/gatling-backend/docs"
	"github.com/mr-isik/gatling-backend/internal/api/handler"
	"github.com/mr-isik/gatling-backend/internal/api/middleware"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *handler.AuthHandler,
	scenarioHandler *handler.ScenarioHandler,
	runHandler *handler.TestRunHandler,
	reportHandler *handler.ReportHandler,
	projectHandler *handler.ProjectHandler,
	wsHandler *handler.WSHandler,
	jwtSecret string,
) {
	// Global Middleware
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Second,
	}))

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Public Routes
	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Protected Routes setup
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// Auth API Keys
	protectedAuth := auth.Group("/", authMiddleware)
	protectedAuth.Post("/api-keys", authHandler.CreateAPIKey)
	protectedAuth.Delete("/api-keys/:id", authHandler.DeleteAPIKey)

	v1 := app.Group("/v1", authMiddleware)

	// Scenarios
	scenarios := v1.Group("/scenarios")
	scenarios.Get("/", scenarioHandler.List)
	scenarios.Post("/", scenarioHandler.Create)
	scenarios.Post("/generate", scenarioHandler.Generate)
	scenarios.Get("/:id", scenarioHandler.Get)
	scenarios.Put("/:id", scenarioHandler.Update)
	scenarios.Delete("/:id", scenarioHandler.Delete)
	scenarios.Post("/:id/clone", scenarioHandler.Clone)

	// Test Runs
	runs := v1.Group("/runs")
	runs.Post("/", runHandler.Start)
	runs.Get("/", runHandler.List)
	runs.Get("/:id", runHandler.Get)
	runs.Post("/:id/stop", runHandler.Stop)
	runs.Get("/:id/metrics", runHandler.Metrics)
	runs.Get("/:id/logs", runHandler.Logs)

	// Reports
	reports := v1.Group("/reports")
	reports.Get("/compare", reportHandler.Compare)
	reports.Get("/:id", reportHandler.Get)
	reports.Get("/:id/ai-summary", reportHandler.AISummary)
	reports.Post("/:id/export", reportHandler.Export)

	// Projects
	projects := v1.Group("/projects")
	projects.Get("/", projectHandler.List)
	projects.Post("/", projectHandler.Create)
	projects.Post("/:id/members", projectHandler.AddMember)
	projects.Get("/:id/usage", projectHandler.Usage)

	// WS Routes (Public or protected via token query param. For now, public stub)
	ws := app.Group("/ws/runs", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	ws.Get("/:id/live", websocket.New(wsHandler.ServeLiveMetrics))
	ws.Get("/:id/anomalies", websocket.New(wsHandler.ServeAnomalies))
}
