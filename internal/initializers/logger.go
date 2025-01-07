package initializers

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func ConfigLogger(ech *echo.Echo) zerolog.Logger {
	var logger zerolog.Logger

	if os.Getenv("APP_ENV") != "development" {
		logger = zerolog.New(os.Stderr).
			With().
			Timestamp().
			Logger()
	} else {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	ech.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRoutePath: true,
		LogStatus:    true,
		LogLatency:   true,
		LogHost:      true,
		LogMethod:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("http_method", v.Method).
				Str("route", v.RoutePath).
				Int("status", v.Status).
				Str("latency", v.Latency.String()).
				Str("host", v.Host).
				Msg("request")

			return nil
		},
	}))

	return logger
}
