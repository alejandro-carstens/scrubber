package controllers

import (
	"errors"

	"github.com/labstack/echo/v4"
)

func registry() []Controllerable {
	return []Controllerable{
		&GoogleAuthController{},
		&UserController{},
		&RoleController{},
		&RolePermissionController{},
		&UserRoleController{},
	}
}

func Bootstrap(e *echo.Echo) error {
	for _, controller := range registry() {
		for _, route := range controller.new().Routes() {
			switch route.Method() {
			case "GET":
				e.GET(route.Route(), route.Handler(), route.Middlewares()...)
				break
			case "POST":
				e.POST(route.Route(), route.Handler(), route.Middlewares()...)
				break
			case "PUT":
				e.PUT(route.Route(), route.Handler(), route.Middlewares()...)
				break
			case "PATCH":
				e.PATCH(route.Route(), route.Handler(), route.Middlewares()...)
				break
			case "DELETE":
				e.DELETE(route.Route(), route.Handler(), route.Middlewares()...)
				break
			default:
				return errors.New("invalid route method specified")
			}
		}
	}

	return nil
}
