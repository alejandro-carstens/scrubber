package app

import (
	"fmt"
	"scrubber/app/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(port int) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if err := controllers.Bootstrap(e); err != nil {
		return err
	}

	return e.Start(fmt.Sprintf(":%v", port))
}
