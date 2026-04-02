package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

type TraceConfig struct {
	Enabled bool
}

func TraceLogger(config TraceConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(contextObj echo.Context) error {
			if !config.Enabled {
				return next(contextObj)
			}

			start := time.Now()
			request := contextObj.Request()
			ctx := request.Context()
			path := request.URL.Path
			method := request.Method

			slog.DebugContext(ctx, "Entry: TraceLogger", "method", method, "path", path)

			err := next(contextObj)

			duration := time.Since(start)
			slog.DebugContext(ctx, "Exit: TraceLogger", "method", method, "path", path, "duration", duration, "status", contextObj.Response().Status)

			return err
		}
	}
}
