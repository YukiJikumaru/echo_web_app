package main

import (
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/YukiJikumaru/echo_web_app/handlers"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

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

func setLogging(e *echo.Echo) *echo.Echo {
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)
	return e
}

func setMddlewares(e *echo.Echo) *echo.Echo {
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
	return e
}

func setRenderer(e *echo.Echo) *echo.Echo {
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

	return e
}

func setRouter(e *echo.Echo) *echo.Echo {
	e.Static("/assets", "public/assets")

	e.GET("/", handlers.GetIndexHandler)
	e.GET("/test", func(c echo.Context) error {
		return c.Render(http.StatusOK, "test.html", echo.Map{"Title": "テストページ"})
	})
	e.GET("/login", handlers.GetLoginHandler)
	e.POST("/login", handlers.PostLoginHandler)
	e.GET("/signup", handlers.GetSignupHandler)
	e.POST("/signup", handlers.PostSignupHandler)
	e.GET("/logout", handlers.GetLogoutHandler)

	return e
}

func Boot(e *echo.Echo) *echo.Echo {
	e = setLogging(e)
	e = setMddlewares(e)
	e = setRenderer(e)
	e = setRouter(e)
	return e
}
