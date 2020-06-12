package contexts

import "errors"

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
