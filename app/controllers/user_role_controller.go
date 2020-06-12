package controllers

import (
	"net/http"
	"scrubber/app/services/accesscontrol"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/labstack/echo/v4"
)

type UserRoleController struct {
	service *accesscontrol.UserRoleService
}

func (rc *UserRoleController) new() Controllerable {
	return &UserRoleController{
		service: accesscontrol.NewUserRoleService(),
	}
}

// Routes implementation of the Controllable interface
func (rc *UserRoleController) Routes() []*Route {
	return []*Route{
		&Route{
			method:  "PUT",
			route:   "/api/users/:user_id/roles/:role_id",
			handler: rc.Handle,
		},
	}
}

func (rc *UserRoleController) Handle(ctx echo.Context) error {
	request := echo.Map{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	context, err := contexts.NewUserRoleContext(request)

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	user, err := rc.service.Handle(context)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"user": user})
}
