package service

import (
	"fmt"
	"github.com/the-psyducks/metrics-service/src/app_errors"
	"github.com/the-psyducks/metrics-service/src/models"
	"github.com/the-psyducks/metrics-service/src/repository"
	"net/http"
	"os"
)

type MetricsService struct {
	database *repository.MetricsPostgresDB
}

func (s MetricsService) RecordLoginAttempt(loginAttempt models.LoginAttempt) *app_errors.AppError {
	err := s.database.RegisterLoginAttempt(loginAttempt)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, fmt.Errorf("error registering login attempt: %w",
			err))
	}
	return nil
}

func (s MetricsService) GetLoginMetrics(isAdmin bool) (*models.LoginSummaryMetrics, error) {

	if !isAdmin && os.Getenv("ENV") == "HEROKU" {
		return nil, app_errors.NewAppError(http.StatusForbidden, app_errors.UserIsNotAdmin, app_errors.ErrUserIsNotAdmin)
	}

	metrics, err := s.database.GetLoginSummaryMetrics()
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, err)
	}
	return metrics, nil
}

func NewMetricsService(db *repository.MetricsPostgresDB) *MetricsService {
	return &MetricsService{db}
}
