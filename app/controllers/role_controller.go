package controllers

import (
	"net/http"
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/accesscontrol"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/labstack/echo/v4"
)

type RoleController struct {
	service    *accesscontrol.RoleService
	repository *repositories.RoleRepository
}

func (rc *RoleController) new() Controllerable {
	return &RoleController{
		service:    accesscontrol.NewRoleService(),
		repository: repositories.NewRoleRepository(),
	}
}

// Routes implementation of the Controllable interface
func (rc *RoleController) Routes() []*Route {
	return []*Route{
		&Route{
			method:  "POST",
			route:   "/api/roles",
			handler: rc.Handle,
		},
		&Route{
			method:  "GET",
			route:   "/api/roles",
			handler: rc.Index,
		},
	}
}

func (rc *RoleController) Handle(ctx echo.Context) error {
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

func (rc *RoleController) Index(ctx echo.Context) error {
	queryContext, err := repositories.BindQueryContext(ctx.QueryParam("query"))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	roles := []*models.Role{}

	meta, err := rc.repository.QueryByContext(queryContext, &roles)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"meta": meta, "roles": roles})
}
