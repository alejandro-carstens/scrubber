package contexts

import (
	"encoding/json"
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

	rc := &RoleContext{
		roleId:      role.RoleId,
		permissions: role.Permissions,
	}

	return rc, rc.validate()
}

func inStringSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}
