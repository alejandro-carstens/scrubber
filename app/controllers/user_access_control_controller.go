package controllers

import (
	"net/http"
	"scrubber/app/services/accesscontrol"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/labstack/echo/v4"
)

type UserAccessControlController struct {
	service *accesscontrol.AccessControlService
}

func (uacc *UserAccessControlController) new() Controllerable {
	return &UserAccessControlController{
		service: accesscontrol.NewAccessControlService(),
	}
}

// Routes implementation of the Controllable interface
func (uacc *UserAccessControlController) Routes() []*Route {
	return []*Route{
		&Route{
			method:  "PUT",
			route:   "/api/users/:user_id/access_controls",
			handler: uacc.Handle,
		},
	}
}

func (uacc *UserAccessControlController) Handle(ctx echo.Context) error {
	request := echo.Map{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	context, err := contexts.NewAccessControlContext(request)

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	user, err := uacc.service.Handle(context)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"user": user})
}
