package handlers

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func LoadRoutes(e *echo.Echo) {
	e.GET("/", homePage)
	e.GET("/login", loginPage)
	e.POST("/login", createSession)
	e.GET("/showup", readSession)
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
