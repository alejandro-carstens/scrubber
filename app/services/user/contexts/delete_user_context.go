package contexts

import (
	"errors"
	"fmt"
	"strconv"
)

func NewDeleteUserContext(params map[string]interface{}) (*DeleteUserContext, error) {
	if _, valid := params["user_id"]; !valid {
		return nil, errors.New("no user_id specified")
	}

	userId, err := strconv.Atoi(fmt.Sprint(params["user_id"]))

	if err != nil {
		return nil, err
	}

	ctx := &DeleteUserContext{userId: uint64(userId)}

	return ctx, ctx.validate()
}

type DeleteUserContext struct {
	userId uint64
}

func (duc *DeleteUserContext) UserID() uint64 {
	return duc.userId
}

func (duc *DeleteUserContext) validate() error {
	if duc.userId <= 0 {
		return errors.New("user_id needs to be greater than 0")
	}

	return nil
}
