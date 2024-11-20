package repository

import (
	//"database/sql"
	//"errors"
	"fmt"

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

	schemaLoginMetrics := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			user_id UUID NOT NULL,
			login_time TIMESTAMPTZ NOT NULL DEFAULT now(),
			succesfull BOOLEAN NOT NULL,
			identity_provider VARCHAR(255) DEFAULT NULL,
			PRIMARY KEY (user_id, login_time),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
			);
		`, "login_metrics")

	if _, err := db.Exec(schemaLoginMetrics); err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}
