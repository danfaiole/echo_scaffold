package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/danfaiole/erp_go/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Log *zerolog.Logger

// LoadRoutes does what the name says, it inserts all the routes into the
// echo instance. We're using structs to organize better the handlers, so please
// keep with the defaults.
func LoadRoutes(e *echo.Echo, dbPool *pgxpool.Pool, logger *zerolog.Logger) {
	Log = logger
	// Load the queries
	queries := database.New(dbPool)
	// Session Routes
	sessHandler := SessionHandler{
		queries: queries,
		ctxBack: context.Background(),
	}
	g := e.Group("/")
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)

			if sess != nil {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}
	})

	g.GET("/", sessHandler.homePage)
	g.GET("/login", sessHandler.loginPage)
	g.POST("/login", sessHandler.createSession)
	g.GET("/showup", sessHandler.readSession)
}

// render is a shortcut function to the render function for templ
func render(ctx echo.Context, statusCode int, component templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := component.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
