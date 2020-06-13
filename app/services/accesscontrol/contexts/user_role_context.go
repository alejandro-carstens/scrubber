package contexts

import (
	"errors"
	"fmt"
	"strconv"
)

func NewUserRoleContext(params map[string]interface{}) (*UserRoleContext, error) {
	if _, valid := params["role_id"]; !valid {
		return nil, errors.New("no role_id specified")
	}

	roleId, err := strconv.Atoi(fmt.Sprint(params["role_id"]))

	if err != nil {
		return nil, err
	}

	if _, valid := params["user_id"]; !valid {
		return nil, errors.New("no user_id specified")
	}

	userId, err := strconv.Atoi(fmt.Sprint(params["user_id"]))

	if err != nil {
		return nil, err
	}

	ctx := &UserRoleContext{
		userId: uint64(userId),
		roleId: uint64(roleId),
	}

	return ctx, ctx.validate()
}

type UserRoleContext struct {
	userId uint64
	roleId uint64
}

func (urc *UserRoleContext) UserID() uint64 {
	return urc.userId
}

func (urc *UserRoleContext) RoleID() uint64 {
	return urc.roleId
}

func (urc *UserRoleContext) validate() error {
	if urc.roleId <= 0 {
		return errors.New("role_id needs to be greater than 0")
	}

	if urc.userId <= 0 {
		return errors.New("user_id needs to be greater than 0")
	}

	return nil
}
