package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echo.Map{"data": "Hello, it works.."})
	})

	// connect to DB
	// configs.ConnectDB()

	e.Logger.Fatal(e.Start("127.0.0.1:8000"))

}
