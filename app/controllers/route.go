package controllers

import "github.com/labstack/echo/v4"

type Route struct {
	method      string
	route       string
	handler     echo.HandlerFunc
	middlewares []echo.MiddlewareFunc
}

func (r *Route) Method() string {
	return r.method
}

func (r *Route) Route() string {
	return r.route
}

func (r *Route) Handler() echo.HandlerFunc {
	return r.handler
}

func (r *Route) Middlewares() []echo.MiddlewareFunc {
	if len(r.middlewares) > 0 {
		return r.middlewares
	}

	return []echo.MiddlewareFunc{}
}
