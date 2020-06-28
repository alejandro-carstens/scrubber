package contexts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func NewUserRoleContext(params map[string]interface{}) (*UserRoleContext, error) {
	if _, valid := params["user_id"]; valid {
		userId, err := strconv.Atoi(fmt.Sprint(params["user_id"]))

		if err != nil {
			return nil, err
		}

		params["user_id"] = uint64(userId)
	} else {
		params["user_id"] = uint64(0)
	}

	b, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	entity := struct {
		UserId  uint64   `json:"user_id"`
		RoleIds []uint64 `json:"role_ids"`
	}{}

	if err := json.Unmarshal(b, &entity); err != nil {
		return nil, err
	}

	ctx := &UserRoleContext{
		userId:  entity.UserId,
		roleIds: entity.RoleIds,
	}

	return ctx, ctx.validate()
}

type UserRoleContext struct {
	userId  uint64
	roleIds []uint64
}

func (urc *UserRoleContext) UserID() uint64 {
	return urc.userId
}

func (urc *UserRoleContext) RoleIDs() []uint64 {
	return urc.roleIds
}

func (urc *UserRoleContext) RoleIDsMap() map[uint64]uint64 {
	roleIdsMap := map[uint64]uint64{}

	for _, roleId := range urc.roleIds {
		roleIdsMap[roleId] = roleId
	}

	return roleIdsMap
}

func (urc *UserRoleContext) validate() error {
	if urc.userId <= 0 {
		return errors.New("user_id needs to be greater than 0")
	}

	return nil
}
