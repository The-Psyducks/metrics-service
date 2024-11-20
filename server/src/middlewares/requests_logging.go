package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"strconv"
	"time"
)

// Middleware for logging requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		formattedDuration := formatDuration(duration)

		logRequest(method, formattedDuration, path, uint(status))
	}
}

func logRequest(method string, time_taken string, endpoint string, status uint) {
	var statusMessage string
	if status >= 200 && status < 300 {
		statusMessage = "Successful request"
	} else {
		statusMessage = "Failed request"
	}

	slog.Info(
		statusMessage,
		slog.String("method", method),
		slog.Int("status", int(status)),
		slog.String("endpoint", endpoint),
		slog.String("time_taken", time_taken),
	)
}

// formats the duration of the request to a round number using ms/µs/ns
func formatDuration(duration time.Duration) (formatted string) {
	micros := duration.Microseconds()
	millis := duration.Milliseconds()
	nanos := duration.Nanoseconds()

	if micros >= 1000 {
		if millis >= 1 {
			return strconv.FormatInt(millis, 10) + " ms"
		}
		return strconv.FormatInt(micros/1000, 10) + " ms"
	}
	if micros >= 1 {
		return strconv.FormatInt(micros, 10) + " µs"
	}
	return strconv.FormatInt(nanos/1000, 10) + " ns"
}