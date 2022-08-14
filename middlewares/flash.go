package middlewares

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

var (
	FlashPath       = "/"
	FlashMaxAge     = 100
	FlashSuccessKey = "flash_success"
	FlashErrorKey   = "flash_error"
	FlashSecretKey  = "secret-session-key"
)

func SetFlash(c echo.Context, name, value string) error {
	session, err := sessions.NewCookieStore([]byte(FlashSecretKey)).Get(c.Request(), name)
	if err != nil {
		return err
	}

	session.AddFlash(value, name)
	session.Options = &sessions.Options{
		Path:     FlashPath,
		MaxAge:   FlashMaxAge,
		HttpOnly: true,
		Secure:   false,
	}

	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}

	return nil
}

func GetFlash(c echo.Context, name string) (flashes []string, err error) {
	session, err := sessions.NewCookieStore([]byte(FlashSecretKey)).Get(c.Request(), name)
	if err != nil {
		return nil, err
	}

	fls := session.Flashes(name)
	if len(fls) > 0 {
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			return nil, err
		}
		for _, fl := range fls {
			flashes = append(flashes, fl.(string))
		}

		return flashes, nil
	}
	return []string{}, nil
}

func SetFlashSuccess(c echo.Context, message string) error {
	return SetFlash(c, FlashSuccessKey, message)
}

func GetFlashSuccess(c echo.Context) ([]string, error) {
	return GetFlash(c, FlashSuccessKey)
}

func SetFlashError(c echo.Context, message string) error {
	return SetFlash(c, FlashErrorKey, message)
}

func GetFlashError(c echo.Context) ([]string, error) {
	return GetFlash(c, FlashErrorKey)
}
