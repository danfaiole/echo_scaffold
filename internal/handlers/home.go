package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/danfaiole/erp_go/internal/database"
	"github.com/danfaiole/erp_go/internal/views/pages"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	queries *database.Queries
	ctxBack context.Context
}

func (h SessionHandler) homePage(ctx echo.Context) error {
	return render(ctx, http.StatusOK, pages.Home())
}

func (h SessionHandler) loginPage(ctx echo.Context) error {
	return render(ctx, http.StatusOK, pages.Login())
}

func (h SessionHandler) createSession(ctx echo.Context) error {
	sess, err := session.Get("session", ctx)

	log.Println(ctx.FormValue("email"))
	log.Println(ctx.FormValue("password1"))

	h.queries.CreateUser(h.ctxBack, database.CreateUserParams{
		Username: ctx.FormValue("email"),
		Password: ctx.FormValue("password1"),
		Email:    ctx.FormValue("email"),
	})

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

	return h.readSession(ctx)
}

func (h SessionHandler) readSession(ctx echo.Context) error {
	sess, err := session.Get("session", ctx)
	if err != nil {
		return err
	}

	resp, err := h.queries.GetUser(h.ctxBack, 1)
	if err != nil {
		logger := ctx.Logger()
		logger.Info(resp)
	}
	log.Println(resp)

	return ctx.String(http.StatusOK, fmt.Sprintf("foo=%v\n", sess.Values["test"]))
}
