package repository

import (
	//"database/sql"
	//"errors"
	"fmt"
	"github.com/the-psyducks/metrics-service/src/config"
	"github.com/the-psyducks/metrics-service/src/models"
	"os"
	"testing"

	//uuid "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	//"github.com/the-psyducks/metrics-service/src/models"
)

type MetricsPostgresDB struct {
	db *sqlx.DB
}

func CreateMetricsPostgresDB(db *sqlx.DB) (*MetricsPostgresDB, error) {
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	postgresDB := MetricsPostgresDB{db}

	return &postgresDB, nil
}

func createTables(db *sqlx.DB) error {

	schemaLoginMetrics := `
		DROP TABLE IF EXISTS login_metrics;

		CREATE TABLE IF NOT EXISTS login_metrics (
			user_id UUID NOT NULL,
			login_time TIMESTAMPTZ NOT NULL DEFAULT now(),
			succesfull BOOLEAN NOT NULL,
			identity_provider VARCHAR(255) DEFAULT NULL,
			PRIMARY KEY (user_id, login_time)
			
			);
		`

	if _, err := db.Exec(schemaLoginMetrics); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

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

func CreateMetricsDatabases(cfg *config.Config) (*MetricsPostgresDB, error) {
	conn, err := createDBConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create user database connection: %w", err)
	}
	db, err := CreateMetricsPostgresDB(conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *MetricsPostgresDB) RegisterLoginAttempt(loginAttempt models.LoginAttempt) error {
	query := `
		INSERT INTO login_metrics (user_id, login_time, succesfull, identity_provider)
		VALUES ($1, NOW(), $2, $3)
	`

	_, err := db.db.Exec(query, loginAttempt.UserId, loginAttempt.WasSuccessful, loginAttempt.Provider)

	if err != nil {
		return fmt.Errorf("failed to register login: %w", err)
	}

	return nil
}

// TODO
func (db *MetricsPostgresDB) GetLoginSummaryMetrics() (*models.LoginSummaryMetrics, error) {
	var loginSummary models.LoginSummaryMetrics

	query := `
		SELECT 
			COUNT(*) AS total_logins,
			COALESCE(SUM(CASE WHEN succesfull THEN 1 ELSE 0 END), 0) AS succesfull_logins,
			COALESCE(SUM(CASE WHEN NOT succesfull THEN 1 ELSE 0 END), 0) AS failed_logins
		FROM login_metrics
	`
	if err := db.db.Get(&loginSummary, query); err != nil {
		return nil, fmt.Errorf("error getting login metrics: %w", err)
	}

	query = `
		SELECT 
			COALESCE(SUM(CASE WHEN identity_provider IS NULL THEN 1 ELSE 0 END), 0) AS email,
			COALESCE(SUM(CASE WHEN identity_provider IS NOT NULL THEN 1 ELSE 0 END), 0) AS federated
		FROM login_metrics
		WHERE succesfull = true
	`
	if err := db.db.Get(&loginSummary.MethodDistribution, query); err != nil {
		return nil, fmt.Errorf("error getting login method distribution: %w", err)
	}

	var federatedProviders []struct {
		Provider string `db:"identity_provider"`
		Amount   int    `db:"amount"`
	}
	query = `
		SELECT identity_provider, COUNT(*) AS amount
		FROM login_metrics
		WHERE identity_provider IS NOT NULL
		GROUP BY identity_provider
	`
	if err := db.db.Select(&federatedProviders, query); err != nil {
		return nil, fmt.Errorf("error getting federated providers: %w", err)
	}

	loginSummary.FederatedProviders = make(map[string]int)
	for _, provider := range federatedProviders {
		loginSummary.FederatedProviders[provider.Provider] = provider.Amount
	}

	return &loginSummary, nil
}
