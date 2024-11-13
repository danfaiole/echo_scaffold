package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danfaiole/erp_go/internal/views/pages"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func homePage(ctx echo.Context) error {
	return render(ctx, http.StatusOK, pages.Home())
}

func loginPage(ctx echo.Context) error {
	return render(ctx, http.StatusOK, pages.Login())
}

func createSession(ctx echo.Context) error {
	sess, err := session.Get("session", ctx)

	log.Println(ctx.FormValue("email"))
	log.Println(ctx.FormValue("password1"))

	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["test"] = "Logged user!"
	if err := sess.Save(ctx.Request(), ctx.Response()); err != nil {
		return err
	}

	return readSession(ctx)
}

func readSession(ctx echo.Context) error {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return err
	}
	return ctx.String(http.StatusOK, fmt.Sprintf("foo=%v\n", sess.Values["test"]))
}
