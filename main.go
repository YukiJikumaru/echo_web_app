package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	app := Boot(echo.New())

	app.Logger.Fatal(app.Start(":1323"))
}
