package controllers

import (
	"net/http"
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/user"
	"scrubber/app/services/user/contexts"

	"github.com/labstack/echo/v4"
)

// UserController controls access and
// edits to the user model
type UserController struct {
	repository        *repositories.UserRepository
	addUserService    *user.AddUserService
	deleteUserService *user.DeleteUserService
}

func (uc *UserController) new() Controllerable {
	return &UserController{
		repository:        repositories.NewUserRepository(),
		addUserService:    user.NewAddUserService(),
		deleteUserService: user.NewDeleteUserService(),
	}
}

// Routes implementation of the Controllable interface
func (uc *UserController) Routes() []*Route {
	return []*Route{
		&Route{
			method:  "POST",
			route:   "/api/users",
			handler: uc.Create,
		},
		&Route{
			method:  "GET",
			route:   "/api/users",
			handler: uc.Index,
		},
		&Route{
			method:  "DELETE",
			route:   "/api/users/:user_id",
			handler: uc.Delete,
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

// Create is responsible for creating a new user
func (uc *UserController) Create(ctx echo.Context) error {
	request := echo.Map{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	context, err := contexts.NewAddUserContext(request)

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	user, err := uc.addUserService.Handle(context)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"user": user})
}

// Delete deletes a given user
func (uc *UserController) Delete(ctx echo.Context) error {
	request := echo.Map{}

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	context, err := contexts.NewDeleteUserContext(request)

	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"error": true, "message": err.Error()})
	}

	if err := uc.deleteUserService.Handle(context); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": true, "message": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"deleted": true})
}
