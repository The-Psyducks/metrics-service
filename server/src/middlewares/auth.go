package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/the-psyducks/metrics-service/src/app_errors"
	"github.com/the-psyducks/metrics-service/src/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Error("Authorization header is required")
			err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("authorization header is required"))
			_ = c.AbortWithError(err.Code, err)
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			slog.Error("Invalid authorization header")
			err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", fmt.Errorf("invalid authorization header"))
			_ = c.AbortWithError(err.Code, err)
			c.Next()
			return
		}

		tokenString := bearerToken[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			slog.Error("Invalid token")
			err := app_errors.NewAppError(http.StatusUnauthorized, "Unauthorized", err)
			_ = c.AbortWithError(err.Code, err)
			return
		}

		c.Set("session_user_id", claims.UserId)
		c.Set("session_user_admin", claims.UserAdmin)
		c.Set("token", tokenString)

		c.Next()
	}
}
