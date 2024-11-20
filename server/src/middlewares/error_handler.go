package middleware

import (
	"errors"
	"github.com/the-psyducks/metrics-service/src/app_errors"
	model "github.com/the-psyducks/metrics-service/src/models"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *app_errors.AppError
			if errors.As(err, &appErr) {
				slog.Error(appErr.Message, slog.String("error", appErr.Error()))
				if appErr.Code == http.StatusInternalServerError {
					sendInternalServerErrorResponse(c)
				} else {
					sendErrorResponse(c, appErr.Code, appErr.Message, appErr.Error())
				}
			}
			c.Abort()
		}
	}
}

func sendErrorResponse(ctx *gin.Context, status int, title string, detail string) {
	errorResponse := model.ErrorResponse{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: ctx.Request.URL.Path,
	}
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(status, errorResponse)
}

func sendInternalServerErrorResponse(ctx *gin.Context) {
	sendErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error", "internal server error - please contact support")
}
