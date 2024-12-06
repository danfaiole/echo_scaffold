package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/danfaiole/erp_go/internal/handlers"
	"github.com/danfaiole/erp_go/internal/initializers"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ech := echo.New()

	// Loads system dependencies
	initializers.LoadEnvVars()
	dbPool := initializers.ConnectDB()

	defer dbPool.Close()

	// Serve static files like js, css
	ech.Static("/static", "assets")

	// Middleware list
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ech.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
	ech.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("ERP_GO_COOKIE_SECRET")))))

	// Loads endpoints into echo instance
	handlers.LoadRoutes(ech, dbPool)

	ech.Logger.Fatal(ech.Start(":1323"))
}
