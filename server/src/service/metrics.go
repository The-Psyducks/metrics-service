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

func NewMetricsService(db *repository.MetricsPostgresDB) *MetricsService {
	return &MetricsService{db}
}

func (s *MetricsService) RecordLoginAttempt(loginAttempt models.LoginAttempt) *app_errors.AppError {
	err := s.database.RegisterLoginAttempt(loginAttempt)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, fmt.Errorf("error registering login attempt: %w",
			err))
	}
	return nil
}

func (s *MetricsService) RecordUserBlocked(userBlocked models.UserBlocked) *app_errors.AppError {
	err := s.database.RegisterUserBlocked(userBlocked)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, fmt.Errorf("error registering blocked user: %w",
			err))
	}
	return nil

}

func (s *MetricsService) RecordUserUnblocked(userUnblocked models.UserUnblocked) *app_errors.AppError {
	err := s.database.RegisterUserUnblocked(userUnblocked)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, fmt.Errorf("error registering unblocked user: %w",
			err))
	}
	return nil

}

func (s *MetricsService) RecordNewRegistry(newRegistry models.NewRegistry) *app_errors.AppError {
	err := s.database.RegisterNewRegistry(newRegistry)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, fmt.Errorf("error registering new registry: %w",
			err))
	}
	return nil
}

func (s *MetricsService) RecordNewUser(newUser models.NewUser) *app_errors.AppError {
	err := s.database.RegisterNewUser(newUser)
	if err != nil {
		return app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, fmt.Errorf("error registering new user: %w",
			err))
	}
	return nil
}

func (s *MetricsService) GetLoginMetrics(isAdmin bool) (*models.LoginSummaryMetrics, error) {

	if !isAdmin && os.Getenv("ENV") == "HEROKU" {
		return nil, app_errors.NewAppError(http.StatusForbidden, app_errors.UserIsNotAdmin, app_errors.ErrUserIsNotAdmin)
	}

	metrics, err := s.database.GetLoginSummaryMetrics()
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, err)
	}
	return metrics, nil
}

func (s *MetricsService) GetRegistryMetrics(admin bool) (*models.RegistrationSummaryMetrics, error) {
	if !admin && os.Getenv("ENV") == "HEROKU" {
		return nil, app_errors.NewAppError(http.StatusForbidden, app_errors.UserIsNotAdmin, app_errors.ErrUserIsNotAdmin)
	}
	metrics, err := s.database.GetRegistrySummaryMetrics()
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, err)
	}
	return metrics, nil

}

func (s *MetricsService) GetLocationMetrics(admin bool) (*models.LocationMetrics, error) {
	if !admin && os.Getenv("ENV") == "HEROKU" {
		return nil, app_errors.NewAppError(http.StatusForbidden, app_errors.UserIsNotAdmin, app_errors.ErrUserIsNotAdmin)
	}
	metrics, err := s.database.GetLocationMetrics()
	if err != nil {
		return nil, app_errors.NewAppError(http.StatusInternalServerError, app_errors.InternalServerError, err)
	}
	return metrics, nil
}
