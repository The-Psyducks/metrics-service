package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/the-psyducks/metrics-service/src/repository"
	"github.com/the-psyducks/metrics-service/src/service"
	"log/slog"
	"net/http"
)

type WebController struct {
	db      *repository.MetricsPostgresDB
	service *service.MetricsService
}

func (c *WebController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}

func (c *WebController) GetLoginMetrics(ctx *gin.Context) {
	userSessionIsAdmin := ctx.GetBool("session_user_admin")
	metrics, err := c.service.GetLoginMetrics(userSessionIsAdmin)
	if err != nil {
		slog.Warn(fmt.Sprintf("error getting login metrics: %v", err))
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, metrics)
}

func NewWebController(db *repository.MetricsPostgresDB) *WebController {
	return &WebController{db: db, service: service.NewMetricsService(db)}
}
