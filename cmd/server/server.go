package main

import (
	"os"

	"github.com/danfaiole/erp_go/internal/handlers"
	"github.com/danfaiole/erp_go/internal/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ech := echo.New()

	// Loads all env vars for the system
	utils.LoadEnvVars()

	// Serve static files like js, css
	ech.Static("/static", "assets")

	// Middleware list
	ech.Use(middleware.Logger())
	ech.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("ERP_GO_COOKIE_SECRET")))))

	// Loads endpoints into echo instance
	handlers.LoadRoutes(ech)

	ech.Logger.Fatal(ech.Start(":1323"))
}
