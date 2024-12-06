package handlers

import (
	"context"

	"github.com/a-h/templ"
	"github.com/danfaiole/erp_go/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// LoadRoutes does what the name says, it inserts all the routes into the
// echo instance. We're using structs to organize better the handlers, so please
// keep with the defaults.
func LoadRoutes(e *echo.Echo, dbPool *pgxpool.Pool) {
	// Load the queries
	queries := database.New(dbPool)
	// Session Routes
	sessHandler := SessionHandler{
		queries: queries,
		ctxBack: context.Background(),
	}
	e.GET("/", sessHandler.homePage)
	e.GET("/login", sessHandler.loginPage)
	e.POST("/login", sessHandler.createSession)
	e.GET("/showup", sessHandler.readSession)
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
