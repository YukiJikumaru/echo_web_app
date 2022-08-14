package main

import (
	"crypto/subtle"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/YukiJikumaru/echo_web_app/handlers"
	"github.com/YukiJikumaru/echo_web_app/middlewares"
)

// VIEW
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// ./viewsにディレクトリを作成した場合に追加が必要
var ViewFilePatterns []string = []string{
	"views/*html",
	"views/layouts/*html",
}

func getViewFilePaths() ([]string, error) {
	result := []string{}

	for _, pattern := range ViewFilePatterns {
		paths, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		result = append(result, paths...)
	}

	return result, nil
}

var (
	CSRFContextKey = "csrf"
	CSRFHiddenName = "csrf_token"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	// LOGGING
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	// MIDDLEWARES
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, latency_human =${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(session.MiddlewareWithConfig(session.Config{
		Skipper: middleware.DefaultSkipper,
		Store:   sessions.NewCookieStore([]byte("secret")),
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		ContextKey:     CSRFContextKey,
		TokenLookup:    "header:X-XSRF-TOKEN,form:" + CSRFHiddenName,
		CookiePath:     "/",
		CookieHTTPOnly: true,
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         3600,
		// ContentSecurityPolicy: "default-src 'self'",
	}))

	// VIEW FILES
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	paths, err := getViewFilePaths()
	if err != nil {
		panic("Failed to get view file paths.")
	}

	println("Loading view templates....")
	for _, x := range paths {
		println(x)
	}

	t := &Template{
		templates: template.Must(template.ParseFiles(paths...)),
	}
	e.Renderer = t

	// ROUTING
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	e.Static("/assets", "public/assets")

	e.GET("/", handlers.IndexHandler)
	e.GET("/test", func(c echo.Context) error {
		return c.Render(http.StatusOK, "test.html", echo.Map{"Title": "テストページ"})
	})
	e.GET("/login", func(c echo.Context) error {
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
	})
	e.POST("/login", func(c echo.Context) error {
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
	})

	e.GET("/signup", func(c echo.Context) error {
		println("GET /signup")

		errorMessages, err := middlewares.GetFlashError(c)
		if err != nil {
			panic("FLAAAAAAAAASH")
		}

		return c.Render(http.StatusOK, "signup.html", echo.Map{"Title": "新規登録画面", "CSRFToken": c.Get("csrf"), "ErrorMessages": errorMessages})
	})

	e.POST("/signup", func(c echo.Context) error {
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
	})

	e.GET("/logout", func(c echo.Context) error {
		err := middlewares.LogOut(c)
		if err != nil {
			panic("FLAAAAAAAAASH")
		}

		return c.Redirect(http.StatusFound, "/")
	})

	e.GET("/login/ok", func(c echo.Context) error {
		s, _ := middlewares.GetLoginSession(c)
		msg := fmt.Sprintf("Congratulations! (ID:%d)\n", s.ID)
		return c.String(http.StatusOK, msg)
	})

	return e
}

func main() {
	router := NewRouter()

	router.Logger.Fatal(router.Start(":1323"))
}
