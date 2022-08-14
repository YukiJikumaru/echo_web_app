package handlers

import (
	"net/http"

	"github.com/YukiJikumaru/echo_web_app/middlewares"
	"github.com/labstack/echo/v4"
)

func GetSignupHandler(c echo.Context) error {
	println("GET /signup")

	errorMessages, err := middlewares.GetFlashError(c)
	if err != nil {
		panic("FLAAAAAAAAASH")
	}

	return c.Render(http.StatusOK, "signup.html", echo.Map{"Title": "新規登録画面", "CSRFToken": c.Get("csrf"), "ErrorMessages": errorMessages})
}

func PostSignupHandler(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordConfirm := c.FormValue("passwordConfirm")

	println("email = ", email)
	println("password = ", password)
	println("passwordConfirm = ", passwordConfirm)

	if email == "" || password == "" || passwordConfirm == "" {
		return c.Render(http.StatusOK, "signup.html", echo.Map{"Title": "新規登録画面", "CSRFToken": c.Get("csrf"), "Email": email, "ErrorMessages": []string{"入力が不足しています"}})
	}

	if password != passwordConfirm {
		return c.Render(http.StatusOK, "signup.html", echo.Map{"Title": "新規登録画面", "CSRFToken": c.Get("csrf"), "Email": email, "ErrorMessages": []string{"パスワードが異なっています"}})
	}

	user := middlewares.DataStore.FindByEmail(email)
	if user != nil {
		return c.Render(http.StatusOK, "signup.html", echo.Map{"Title": "新規登録画面", "CSRFToken": c.Get("csrf"), "Email": email, "ErrorMessages": []string{"このメールアドレスは既に登録されています"}})
	}

	middlewares.DataStore.Save(email, password, nil)

	middlewares.SetFlashSuccess(c, "新規登録に成功しました")
	return c.Redirect(http.StatusFound, "/login")
}
