package contexts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"scrubber/app/repositories"
)

func NewAccessControlContext(params map[string]interface{}) (*AccessControlContext, error) {
	if _, valid := params["user_id"]; valid {
		userId, err := strconv.Atoi(fmt.Sprint(params["user_id"]))

		if err != nil {
			return nil, err
		}

		params["user_id"] = uint64(userId)
	}

	b, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	accessControl := struct {
		UserId         uint64           `json:"user_id"`
		AccessControls []*AccessControl `json:"access_controls"`
	}{}

	if err := json.Unmarshal(b, &accessControl); err != nil {
		return nil, err
	}

	acc := &AccessControlContext{
		userId:         accessControl.UserId,
		accessControls: accessControl.AccessControls,
	}

	return acc, acc.validate()
}

type AccessControlContext struct {
	userId         uint64
	accessControls []*AccessControl
}

func (acc *AccessControlContext) UserID() uint64 {
	return acc.userId
}

func (acc *AccessControlContext) AccessControls() []*AccessControl {
	return acc.accessControls
}

func (acc *AccessControlContext) AccessControlsMap() map[string]*AccessControl {
	accessControlsMap := map[string]*AccessControl{}

	for _, accessControl := range acc.accessControls {
		accessControlsMap[accessControl.Action] = accessControl
	}

	return accessControlsMap
}

func (acc *AccessControlContext) validate() error {
	if acc.userId <= 0 {
		return errors.New("invalid user_id, please specify a value greater than 0")
	}

	for _, accessControl := range acc.accessControls {
		if err := accessControl.validate(); err != nil {
			return err
		}
	}

	return nil
}

type AccessControl struct {
	Action string `json:"action"`
	Scope  string `json:"scope"`
}

func (ac *AccessControl) validate() error {
	if !inStringSlice(ac.Action, availableActions) {
		return errors.New(fmt.Sprintf("%v is not a valid action", ac.Action))
	}

	if !inStringSlice(ac.Scope, availableScopes) {
		return errors.New(fmt.Sprintf("%v is not a valid scope", ac.Scope))
	}

	if ac.Action == repositories.ACCESS_CONTROL_ALL_ACTIONS && !inStringSlice(ac.Scope, availableAllActionScopes) {
		return errors.New(fmt.Sprintf("%v is not a valid scope for modifying admin permission", ac.Scope))
	}

	return nil
}
