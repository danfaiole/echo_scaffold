package initializers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var Logger *zerolog.Logger

func ConfigErrors(ech *echo.Echo, logger *zerolog.Logger) {
	Logger = logger
	ech.HTTPErrorHandler = handleErrors
}

func handleErrors(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	Logger.Log().Msg("passed here")
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}
