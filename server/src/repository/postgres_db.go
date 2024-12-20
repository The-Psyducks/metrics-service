package repository

import (
	//"database/sql"
	//"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/the-psyducks/metrics-service/src/config"
	"github.com/the-psyducks/metrics-service/src/models"
	"os"
	"testing"
	"time"

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
		

		CREATE TABLE IF NOT EXISTS login_metrics (
			user_id UUID NOT NULL,
			login_time TIMESTAMPTZ NOT NULL DEFAULT now(),
			succesfull BOOLEAN NOT NULL,
			identity_provider VARCHAR(255) DEFAULT NULL,
			PRIMARY KEY (user_id, login_time)
			);
		`

	schemaUsersBlocked := `
		

		CREATE TABLE IF NOT EXISTS metrics_users_blocks (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL,
			reason TEXT DEFAULT NULL,
			blocked_at TIMESTAMPTZ,
			unblocked_at TIMESTAMPTZ DEFAULT NULL,
			CONSTRAINT users_blocks_unique_block UNIQUE (user_id, blocked_at)
			);
		`

	schemaRegistries := `
		

		CREATE TABLE IF NOT EXISTS metrics_registries (
			registration_id VARCHAR(255) PRIMARY KEY,
			created_at TIMESTAMPTZ,
			deleted_at TIMESTAMPTZ default NULL,			
			identity_provider VARCHAR(255) DEFAULT NULL
			
			);
		`

	schemaUsers := `
					CREATE TABLE IF NOT EXISTS metrics_locations (
					    						user_id UUID PRIMARY KEY,
					    						location VARCHAR(255),
					    						created_at TIMESTAMPTZ)`

	if _, err := db.Exec(schemaLoginMetrics); err != nil {
		return fmt.Errorf("failed to create logins table: %w", err)
	}
	if _, err := db.Exec(schemaUsersBlocked); err != nil {
		return fmt.Errorf("failed to create blocks table: %w", err)
	}
	if _, err := db.Exec(schemaRegistries); err != nil {
		return fmt.Errorf("failed to create registries table: %w", err)
	}

	if _, err := db.Exec(schemaUsers); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
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
		INSERT INTO metrics_login_metrics (user_id, login_time, succesfull, identity_provider)
		VALUES ($1, $2, $3, $4)
	`
	timestamp, err := time.Parse(time.RFC3339, loginAttempt.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %w", err)
	}

	_, err = db.db.Exec(query, loginAttempt.UserId, timestamp, loginAttempt.WasSuccessful, loginAttempt.Provider)

	if err != nil {
		return fmt.Errorf("failed to register login: %w", err)
	}

	return nil
}

func (db *MetricsPostgresDB) RegisterUserBlocked(userBlocked models.UserBlocked) error {
	query := `INSERT INTO metrics_users_blocks (user_id, blocked_at, reason) VALUES ($1, $2, $3)`
	_, err := db.db.Exec(query, userBlocked.UserId, userBlocked.Timestamp, userBlocked.Reason)
	if err != nil {
		return fmt.Errorf("error blocking user: %w", err)
	}
	return nil
}

func (db *MetricsPostgresDB) RegisterUserUnblocked(userUnblocked models.UserUnblocked) error {
	query := `UPDATE metrics_users_blocks SET unblocked_at = $1 
              WHERE user_id = $2 AND unblocked_at IS NULL`
	_, err := db.db.Exec(query, userUnblocked.Timestamp, userUnblocked.UserId)
	if err != nil {
		return fmt.Errorf("error unblocking user: %w", err)
	}
	return nil
}

func (db *MetricsPostgresDB) RegisterNewRegistry(registry models.NewRegistry) error {

	_, err := db.db.Exec("INSERT INTO metrics_registries (registration_id, created_at, identity_provider) VALUES ($1, $2, $3)", registry.RegistrationId, registry.TimeStamp, registry.Provider)
	if err != nil {
		return fmt.Errorf("failed to create registry entry: %w", err)
	}
	return nil
}

func (db *MetricsPostgresDB) RegisterNewUser(newUser models.NewUser) error {
	query := `UPDATE metrics_registries SET deleted_at  = $1
							  WHERE registration_id = $2`

	_, err := db.db.Exec(query, newUser.TimeStamp, newUser.RegistrationId)
	if err != nil {
		return fmt.Errorf("failed to create registry entry: %w", err)
	}

	query = `INSERT INTO metrics_locations (user_id, location, created_at)
    			 VALUES ($1, $2, $3)`

	_, err = db.db.Exec(query, newUser.UserId, newUser.Location, newUser.TimeStamp)
	if err != nil {
		return fmt.Errorf("failed to create registry entry: %w", err)
	}
	return nil
}

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

func (db *MetricsPostgresDB) GetRegistrySummaryMetrics() (*models.RegistrationSummaryMetrics, error) {
	var metrics models.RegistrationSummaryMetrics

	query := `SELECT COUNT(*) as total_registrations,
				COALESCE(SUM(CASE WHEN deleted_at IS NOT NULL THEN 1 ELSE 0 END), 0) as succesfull_registrations,
				COALESCE(SUM(CASE WHEN deleted_at IS NOT NULL THEN 0 ELSE 1 END), 0) as failed_registrations,
				COALESCE(AVG(EXTRACT(EPOCH FROM (deleted_at - created_at))), 0) as average_registration_time
			FROM metrics_registries
	`
	err := db.db.Get(&metrics, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get registration summary metrics: %w", err)
	}

	query = `SELECT
				COALESCE(SUM(CASE WHEN identity_provider = 'INTERNAL' THEN 1 ELSE 0 END), 0) AS email,
				COALESCE(SUM(CASE WHEN identity_provider = 'GOOGLE' THEN 1 ELSE 0 END), 0) AS federated
			FROM metrics_registries
	`
	err = db.db.Get(&metrics.MethodDistribution, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get registration method distribution: %w", err)
	}

	var federatedProviders []struct {
		Provider string `db:"identity_provider"`
		Amount   int    `db:"amount"`
	}
	query = `
		SELECT identity_provider, COUNT(*) AS amount
		FROM metrics_registries
		WHERE identity_provider = 'GOOGLE'
		GROUP BY identity_provider
	`
	if err := db.db.Select(&federatedProviders, query); err != nil {
		return nil, fmt.Errorf("error getting federated providers: %w", err)
	}

	// Inicializa el mapa de federated providers en caso de que esté vacío
	metrics.FederatedProviders = make(map[string]int)
	for _, provider := range federatedProviders {
		metrics.FederatedProviders[provider.Provider] = provider.Amount
	}

	return &metrics, nil
}

func (db *MetricsPostgresDB) GetLocationMetrics() (*models.LocationMetrics, error) {
	var locationMetrics models.LocationMetrics
	query := `
		SELECT location AS country, COUNT(*) AS amount
		FROM metrics_locations
		GROUP BY location
	`

	if err := db.db.Select(&locationMetrics.Locations, query); err != nil {
		return nil, fmt.Errorf("error getting location metrics: %w", err)
	}
	if locationMetrics.Locations == nil {
		locationMetrics.Locations = []models.LocationMetric{}
	}
	return &locationMetrics, nil
}

func (db *MetricsPostgresDB) GetBlockedMetrics() (*models.UsersBlockedMetrics, error) {
	var usersBlockedMetrics models.UsersBlockedMetrics
	query := `
        SELECT 
            COUNT(*) AS total_users_blocked,
            COALESCE(SUM(CASE WHEN unblocked_at IS NULL THEN 1 ELSE 0 END), 0) AS currently_blocked,
            COALESCE(AVG(EXTRACT(EPOCH FROM unblocked_at - blocked_at) / 86400), 0) AS average_block_time_in_days
        FROM metrics_users_blocks
    `

	if err := db.db.Get(&usersBlockedMetrics, query); err != nil {
		return nil, fmt.Errorf("error getting users blocked metrics: %w", err)
	}

	var reasons pq.StringArray
	query = `SELECT ARRAY(SELECT DISTINCT reason FROM metrics_users_blocks WHERE reason IS NOT NULL)`

	if err := db.db.Get(&reasons, query); err != nil {
		return nil, fmt.Errorf("error getting reasons: %w", err)
	}

	usersBlockedMetrics.Reasons = reasons

	return &usersBlockedMetrics, nil
}
