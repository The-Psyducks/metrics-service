package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/the-psyducks/metrics-service/src/repository"
)

type WebController struct {
	db *repository.MetricsPostgresDB
}

func NewWebController(db *repository.MetricsPostgresDB) *WebController {
	return &WebController{db: db}
}

func HealthCheck(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}

func GetMetrics(ctx *gin.Context) {

}
