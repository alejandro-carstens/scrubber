package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var app *App

type App struct {
	e *echo.Echo
}

func Bootstrap() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return nil
}
