package contexts

import (
	"errors"
	"fmt"

	"scrubber/app/repositories"
)

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

func (rc *RoleContext) PermissionMap() map[string]*PermissionEntity {
	permissionMap := map[string]*PermissionEntity{}

	for _, permission := range rc.permissions {
		permissionMap[permission.Action] = permission
	}

	return permissionMap
}

func (rc *RoleContext) validate() error {
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
