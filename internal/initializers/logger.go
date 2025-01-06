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
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
		LogLatency: true,
		LogHost:    true,
	}))

	return logger
}
