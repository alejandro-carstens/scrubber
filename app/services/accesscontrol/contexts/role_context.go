package contexts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"scrubber/app/repositories"
)

func NewRoleContext(params map[string]interface{}) (*RoleContext, error) {
	if _, valid := params["role_id"]; valid {
		roleId, err := strconv.Atoi(fmt.Sprint(params["role_id"]))

		if err != nil {
			return nil, err
		}

		params["role_id"] = uint64(roleId)
	} else {
		params["role_id"] = uint64(0)
	}

	b, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	role := struct {
		RoleId      uint64              `json:"role_id"`
		Name        string              `json:"name"`
		Permissions []*PermissionEntity `json:"permissions"`
	}{}

	if err := json.Unmarshal(b, &role); err != nil {
		return nil, err
	}

	ctx := &RoleContext{
		roleId:      role.RoleId,
		permissions: role.Permissions,
		name:        role.Name,
	}

	return ctx, ctx.validate()
}

type RoleContext struct {
	roleId      uint64
	permissions []*PermissionEntity
	name        string
}

func (rc *RoleContext) RoleID() uint64 {
	return rc.roleId
}

func (rc *RoleContext) Permissions() []*PermissionEntity {
	return rc.permissions
}

func (rc *RoleContext) Name() string {
	return rc.name
}

func (rc *RoleContext) PermissionsMap() map[string]*PermissionEntity {
	permissionsMap := map[string]*PermissionEntity{}

	for _, permission := range rc.permissions {
		permissionsMap[permission.Action] = permission
	}

	return permissionsMap
}

func (rc *RoleContext) validate() error {
	if len(rc.name) == 0 {
		return errors.New("role name cannot be empty")
	}

	if rc.roleId < 0 {
		return errors.New("invalid role_id, please specify a value greater than or equal to 0")
	}

	for _, permission := range rc.permissions {
		if err := permission.validate(); err != nil {
			return err
		}
	}

	return nil
}

type PermissionEntity struct {
	Action string `json:"action"`
	Scope  string `json:"scope"`
}

func (pe *PermissionEntity) validate() error {
	if !inStringSlice(pe.Action, repositories.AvailableActions) {
		return errors.New(fmt.Sprintf("%v is not a valid action", pe.Action))
	}

	if !inStringSlice(pe.Scope, repositories.AvailableScopes) {
		return errors.New(fmt.Sprintf("%v is not a valid scope", pe.Scope))
	}

	return nil
}
