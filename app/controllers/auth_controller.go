package controllers

import (
	"scrubber/app/services/auth"
	"scrubber/app/services/auth/contexts"

	"github.com/labstack/echo/v4"
)

type GoogleAuthController struct {
	service *auth.GoogleAuthService
}

func (gac *GoogleAuthController) new() *GoogleAuthController {
	return &GoogleAuthController{
		service: auth.NewGoogleAuthService(),
	}
}

func (gac *GoogleAuthController) Handle(ctx echo.Context) error {
	request := echo.Map{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(422, echo.Map{"error": true, "message": err.Error()})
	}

	context, err := contexts.NewGoogleAuthContext(request)

	if err != nil {
		return ctx.JSON(422, echo.Map{"error": true, "message": err.Error()})
	}

	token, err := gac.service.Handle(context)

	if err != nil {
		return ctx.JSON(500, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(200, echo.Map{"token": token})
}
