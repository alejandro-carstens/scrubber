package controllers

import (
	"net/http"
	"scrubber/app/services/accesscontrol"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/labstack/echo/v4"
)

type RolePermissionController struct {
	service *accesscontrol.RoleService
}

func (rc *RolePermissionController) new() Controllerable {
	return &RolePermissionController{
		service: accesscontrol.NewRoleService(),
	}
}

// Routes implementation of the Controllable interface
func (rc *RolePermissionController) Routes() []*Route {
	return []*Route{
		&Route{
			method:  "PUT",
			route:   "/api/roles/:role_id/permissions",
			handler: rc.Handle,
		},
	}
}

func (rc *RolePermissionController) Handle(ctx echo.Context) error {
	request := echo.Map{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	context, err := contexts.NewRoleContext(request)

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	role, err := rc.service.Handle(context)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"role": role})
}
