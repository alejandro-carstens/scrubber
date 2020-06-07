package controllers

import (
	"net/http"
	"scrubber/app/models"
	"scrubber/app/repositories"

	"github.com/labstack/echo/v4"
)

// UserController controls access and
// edits to the user model
type UserController struct {
	repository *repositories.UserRepository
}

func (uc *UserController) new() Controllerable {
	return &UserController{
		repository: repositories.NewUserRepository(),
	}
}

// Routes implementation of the Controllable interface
func (uc *UserController) Routes() []*Route {
	return []*Route{
		&Route{
			method:  "GET",
			route:   "/api/users",
			handler: uc.Index,
		},
	}
}

// Index is responsible for returning all
// users for a specified query context
func (uc *UserController) Index(ctx echo.Context) error {
	queryContext, err := repositories.BindQueryContext(ctx.QueryParam("query"))

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	users := []*models.User{}

	meta, err := uc.repository.QueryByContext(queryContext, &users)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"meta": meta, "users": users})
}
