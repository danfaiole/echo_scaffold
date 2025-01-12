package main

import (
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
	logger := initializers.ConfigLogger(ech)
	dbPool := initializers.ConnectDB(&logger)
	initializers.ConfigErrors(ech, &logger)

	defer dbPool.Close()

	// Serve static files like js, css
	ech.Static("/static", "assets")

	// Middleware list
	ech.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("ERP_GO_COOKIE_SECRET")))))
	ech.Use(middleware.Secure())

	// Loads endpoints into echo instance
	handlers.LoadRoutes(ech, dbPool, &logger)

	ech.Logger.Fatal(ech.Start(":1323"))
}
