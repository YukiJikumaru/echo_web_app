package handlers

import (
	"net/http"

	"github.com/YukiJikumaru/echo_web_app/middlewares"
	"github.com/labstack/echo/v4"
)

type IndexLoggedInParams struct {
	Title           string
	SuccessMessages []string
	ErrorMessages   []string
	CurrentUser     *middlewares.User
	MoreStyles      []string
	MoreScripts     []string
}

const Template = "index.html"
const Title = "TOPページ"

func GetIndexHandler(c echo.Context) error {
	println("IsLoggedIn? ", middlewares.IsLoggedIn(c))

	successMessages, _ := middlewares.GetFlashSuccess(c)
	errorMessages, _ := middlewares.GetFlashError(c)
	println("println(successMessages)")
	println(successMessages)
	println("println(errorMessages)")
	println(errorMessages)

	session, err := middlewares.GetLoginSession(c)
	if err == nil {
		println("Session ID = ", session.ID)
	}

	if err != nil {
		return c.Render(http.StatusOK, Template, echo.Map{"Title": Title, "SuccessMessages": successMessages, "ErrorMessages": errorMessages})
	} else {
		user := middlewares.DataStore.FindByID(session.ID)
		println("user = ", user)
		params := IndexLoggedInParams{Title: Title, SuccessMessages: successMessages, ErrorMessages: errorMessages, CurrentUser: user}
		return c.Render(http.StatusOK, Template, params)
	}
}
