package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mr-isik/gatling-backend/internal/api"
	"github.com/mr-isik/gatling-backend/internal/api/handler"
	"github.com/mr-isik/gatling-backend/internal/api/ws"
	"github.com/mr-isik/gatling-backend/internal/engine"
	"github.com/mr-isik/gatling-backend/internal/infra"
	"github.com/mr-isik/gatling-backend/internal/repository/influx"
	"github.com/mr-isik/gatling-backend/internal/repository/postgres"
	"github.com/mr-isik/gatling-backend/internal/repository/redis"
	"github.com/mr-isik/gatling-backend/internal/service"
)

// @title           Gatling-Backend API
// @version         1.0
// @description     AI-Powered Load Testing Tool Backend API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	fmt.Println("Starting Gatling-Backend...")

	// 1. Load Config
	cfg, err := infra.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Init Infra (DB, Redis, Influx)
	db, err := infra.NewDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	rdb, err := infra.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	ifxClient, writeAPI, queryAPI, err := infra.NewInflux(cfg.Influx)
	if err != nil {
		log.Fatalf("Failed to connect to InfluxDB: %v", err)
	}
	defer ifxClient.Close()

	// 3. Init Repos
	scenarioRepo := postgres.NewScenarioRepository(db)
	runRepo := postgres.NewTestRunRepository(db)
	userRepo := postgres.NewUserRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	reportRepo := postgres.NewReportRepository(db)
	metricRepo := influx.NewMetricRepository(writeAPI, queryAPI, cfg.Influx.Bucket)
	cacheRepo := redis.NewCacheRepository(rdb)
	runStateRepo := redis.NewRunStateRepository(rdb)

	// 4. Init Services (Phase 3)
	llmClient := infra.NewLLMClient(cfg.LLM.APIKey)
	aiService := service.NewAIService(llmClient)

	authService := service.NewAuthService(userRepo, cfg.JWT)
	scenarioService := service.NewScenarioService(scenarioRepo, aiService)
	testRunService := service.NewTestRunService(runRepo, scenarioRepo, metricRepo, runStateRepo, cacheRepo)
	reportService := service.NewReportService(reportRepo, metricRepo, aiService)
	baselineService := service.NewBaselineService(runRepo, metricRepo)

	// 4.5 Init Engine (Phase 5) & WS Hub (Phase 6)
	hub := ws.NewHub()
	go hub.Run()

	orchestrator := engine.NewOrchestrator(testRunService, metricRepo, scenarioRepo, aiService, hub)
	testRunService.SetOrchestrator(orchestrator)

	// 5. Init Handlers (Phase 4)
	authHandler := handler.NewAuthHandler(authService)
	scenarioHandler := handler.NewScenarioHandler(scenarioService)
	runHandler := handler.NewTestRunHandler(testRunService)
	reportHandler := handler.NewReportHandler(reportService, baselineService)
	projectHandler := handler.NewProjectHandler(projectRepo)
	wsHandler := handler.NewWSHandler(hub)

	router := api.SetupRoutes(
		authHandler, scenarioHandler, runHandler,
		reportHandler, projectHandler, wsHandler, cfg.JWT.Secret,
	)

	serverPort := cfg.Server.Port
	if serverPort == "" {
		serverPort = "8080"
	}
	log.Printf("Server listening on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
