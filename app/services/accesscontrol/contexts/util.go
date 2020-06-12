package contexts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func NewRoleContext(params map[string]interface{}) (*RoleContext, error) {
	if _, valid := params["role_id"]; valid {
		userId, err := strconv.Atoi(fmt.Sprint(params["user_id"]))

		if err != nil {
			return nil, err
		}

		params["role_id"] = uint64(userId)
	} else {
		params["role_id"] = uint64(0)
	}

	b, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	role := struct {
		RoleId      uint64              `json:"role_id"`
		Permissions []*PermissionEntity `json:"permissions"`
	}{}

	if err := json.Unmarshal(b, &role); err != nil {
		return nil, err
	}

	ctx := &RoleContext{
		roleId:      role.RoleId,
		permissions: role.Permissions,
	}

	return ctx, ctx.validate()
}

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

func inStringSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}
