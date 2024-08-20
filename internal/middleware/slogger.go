package middleware

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

func SLogger(log *slog.Logger, skipPath []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		// process request
		c.Next()

		timestamp := time.Now()
		latency := timestamp.Sub(startTime)
		status := c.Writer.Status()

		slogAttrs := []slog.Attr{
			slog.String("client_ip", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.Int("status_code", status),
			slog.Int("body_size", c.Writer.Size()),
			slog.String("path", path),
			slog.String("query", query),
			slog.Any("params", params),
			slog.Int64("latency(ms)", latency.Milliseconds()),
		}

		level := slog.LevelInfo
		msg := "Incoming request"
		if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			level = slog.LevelWarn
			msg = "client error"
		} else if status >= http.StatusInternalServerError {
			level = slog.LevelError
			msg = "server error"
		}
		if len(c.Errors) > 0 {
			slogAttrs = append(slogAttrs, slog.Any("errors_stack", c.Errors))
		}
		if !slices.Contains(skipPath, path) || status != http.StatusOK {
			log.LogAttrs(c.Request.Context(), level, msg, slogAttrs...)
		}

	}
}
