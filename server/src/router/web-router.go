package router

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/the-psyducks/metrics-service/src/config"
	"github.com/the-psyducks/metrics-service/src/controller"
	middleware "github.com/the-psyducks/metrics-service/src/middlewares"
	"github.com/the-psyducks/metrics-service/src/repository"
	"io"
	"log/slog"
	"os"
)

// Router is a wrapper for the gin.Engine and the address where it is running
type Router struct {
	Engine  *gin.Engine
	Address string
}

// Creates a new router with the configuration provided in the env file
func createRouterFromConfig(cfg *config.Config) *Router {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	router := &Router{
		Engine:  gin.Default(),
		Address: cfg.Host + ":" + cfg.Port,
	}

	return router
}

// Creates a new database connection using the configuration provided in the env file

func addCorsConfiguration(r *Router) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	r.Engine.Use(cors.New(corsConfig))
}

// CreateRouter Creates a new router with the configuration provided in the env file
func CreateRouter() (*Router, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	r := createRouterFromConfig(cfg)

	metricsDb, err := repository.CreateMetricsDatabases(cfg)

	if err != nil {
		slog.Error("failed to create databases", slog.String("error", err.Error()))
		return nil, err
	}

	webController := controller.NewWebController(metricsDb)

	r.Engine.Use(middleware.RequestLogger())
	r.Engine.Use(middleware.ErrorHandler())

	addCorsConfiguration(r)
	r.Engine.GET("/health-check", webController.HealthCheck)

	private := r.Engine.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/metrics/login", webController.GetLoginMetrics)
		private.GET("/metrics/registry", webController.GetRegistryMetrics)
		private.GET("/metrics/location", webController.GetLocationMetrics)
		private.GET("/metrics/blocked", webController.GetBlockedMetrics)
	}

	return r, nil
}

// Run Runs the router in the address provided in the env file
func (r *Router) Run() error {
	fmt.Println("Running in address: ", r.Address)
	return r.Engine.Run(r.Address)
}
