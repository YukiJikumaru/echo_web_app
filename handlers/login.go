package handlers

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/YukiJikumaru/echo_web_app/middlewares"
	"github.com/labstack/echo/v4"
)

func GetLoginHandler(c echo.Context) error {
	println("GET /login")

	errorMessages, err := middlewares.GetFlashError(c)
	if err != nil {
		panic("FLAAAAAAAAASH")
	}
	successMessages, err := middlewares.GetFlashSuccess(c)
	if err != nil {
		panic("FLAAAAAAAAASH")
	}

	return c.Render(http.StatusOK, "login.html", echo.Map{"Title": "ログイン画面", "CSRFToken": c.Get("csrf"), "ErrorMessages": errorMessages, "SuccessMessages": successMessages})
}

func PostLoginHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	println("email = ", email)
	println("password = ", password)

	user := middlewares.DataStore.FindByEmail(email)
	if user == nil {
		middlewares.SetFlashError(c, "ユーザ名かパスワードが間違っています 1")
		return c.Redirect(http.StatusFound, "/login")
	}
	if subtle.ConstantTimeCompare([]byte(password), []byte(user.Pass)) == 1 {
		err := middlewares.SetLoginSession(c, middlewares.LoginSession{ID: user.ID})
		if err != nil {
			fmt.Println("ERROR!!!!!!!", err, "wowow")
		}
		middlewares.SetFlashSuccess(c, "ログインに成功しました")
		return c.Redirect(http.StatusFound, "/")
	} else {
		middlewares.SetFlashError(c, "ユーザ名かパスワードが間違っています 2")
		return c.Redirect(http.StatusFound, "/login")
	}
}
