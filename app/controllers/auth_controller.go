package controllers

import (
	"net/http"

	"scrubber/app/services/auth"
	"scrubber/app/services/auth/contexts"

	"github.com/labstack/echo/v4"
)

// GoogleAuthController represents the
// controller for authenticating
// via Google
type GoogleAuthController struct {
	service *auth.GoogleAuthService
}

func (gac *GoogleAuthController) new() Controllerable {
	return &GoogleAuthController{
		service: auth.NewGoogleAuthService(),
	}
}

// Routes implementation of the Controllable interface
func (gac *GoogleAuthController) Routes() []*Route {
	return []*Route{
		&Route{
			method:  "POST",
			route:   "/api/oauth/google",
			handler: gac.Handle,
		},
	}
}

// Handle is responsible for processing the Google auth process
func (gac *GoogleAuthController) Handle(ctx echo.Context) error {
	request := echo.Map{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	context, err := contexts.NewGoogleAuthContext(request)

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	token, err := gac.service.Handle(context)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"token": token})
}
