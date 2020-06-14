package contexts

import (
	"errors"
	"fmt"
	"strconv"
)

func NewDeleteRoleContext(params map[string]interface{}) (*DeleteRoleContext, error) {
	if _, valid := params["role_id"]; !valid {
		return nil, errors.New("no role_id specified")
	}

	roleId, err := strconv.Atoi(fmt.Sprint(params["role_id"]))

	if err != nil {
		return nil, err
	}

	ctx := &DeleteRoleContext{roleId: uint64(roleId)}

	return ctx, ctx.validate()
}

type DeleteRoleContext struct {
	roleId uint64
}

func (drc *DeleteRoleContext) RoleID() uint64 {
	return drc.roleId
}

func (drc *DeleteRoleContext) validate() error {
	if drc.roleId <= 0 {
		return errors.New("role_id needs to be greater than 0")
	}

	return nil
}
