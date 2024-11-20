package router

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sqlx "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/the-psyducks/metrics-service/src/config"
	"github.com/the-psyducks/metrics-service/src/controller"
	middleware "github.com/the-psyducks/metrics-service/src/middlewares"
	"github.com/the-psyducks/metrics-service/src/repository"
	"io"
	"log/slog"
	"os"
	"testing"
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

func createDBConnection(cfg *config.Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	if testing.Testing() {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DatabaseUser,
			cfg.DatabasePassword,
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseName)

		db, err = sqlx.Connect("postgres", dsn)
	} else {
		switch cfg.Environment {
		case "HEROKU":
			fallthrough
		case "production":
			db, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		case "development":
			fallthrough
		case "testing":
			dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
				cfg.DatabaseUser,
				cfg.DatabasePassword,
				cfg.DatabaseHost,
				cfg.DatabasePort,
				cfg.DatabaseName)

			db, err = sqlx.Connect("postgres", dsn)
		default:
			return nil, fmt.Errorf("invalid environment: %s", cfg.Environment)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	enableUUIDExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := db.Exec(enableUUIDExtension); err != nil {
		return nil, fmt.Errorf("failed to enable uuid extension: %w", err)
	}

	return db, nil
}

func createDatabases(cfg *config.Config) (*repository.MetricsPostgresDB, error) {
	conn, err := createDBConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create user database connection: %w", err)
	}
	db, err := repository.CreateMetricsPostgresDB(conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

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

	metricsDb, err := createDatabases(cfg)

	if err != nil {
		slog.Error("failed to create databases", slog.String("error", err.Error()))
		return nil, err
	}

	controller.NewWebController(metricsDb)

	r.Engine.Use(middleware.RequestLogger())
	r.Engine.Use(middleware.ErrorHandler())

	addCorsConfiguration(r)
	r.Engine.GET("/health-check", controller.HealthCheck)

	private := r.Engine.Group("/")
	private.Use(middleware.AuthMiddleware())
	{
		private.GET("/metrics", controller.GetMetrics)
	}

	return r, nil
}

// Run Runs the router in the address provided in the env file
func (r *Router) Run() error {
	fmt.Println("Running in address: ", r.Address)
	return r.Engine.Run(r.Address)
}
