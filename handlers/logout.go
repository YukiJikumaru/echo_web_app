package handlers

import (
	"net/http"

	"github.com/YukiJikumaru/echo_web_app/middlewares"
	"github.com/labstack/echo/v4"
)

func GetLogoutHandler(c echo.Context) error {
	err := middlewares.LogOut(c)
	if err != nil {
		panic("FLAAAAAAAAASH")
	}

	return c.Redirect(http.StatusFound, "/")
}
