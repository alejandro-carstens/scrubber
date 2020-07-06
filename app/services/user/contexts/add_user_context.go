package contexts

import (
	"errors"
	"regexp"
)

func NewAddUserContext(params map[string]interface{}) (*AddUserContext, error) {
	ctx := &AddUserContext{
		email:    params["email"].(string),
		name:     params["name"].(string),
		lastName: params["last_name"].(string),
	}

	return ctx, ctx.validate()
}

type AddUserContext struct {
	email    string
	name     string
	lastName string
}

func (auc *AddUserContext) Email() string {
	return auc.email
}

func (auc *AddUserContext) Name() string {
	return auc.name
}

func (auc *AddUserContext) LastName() string {
	return auc.lastName
}

func (auc *AddUserContext) validate() error {
	if len(auc.name) == 0 {
		return errors.New("name cannot be empty")
	}

	if len(auc.lastName) == 0 {
		return errors.New("last_name cannot be empty")
	}

	if len(auc.email) == 0 {
		return errors.New("email cannot be empty")
	}

	rgx := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)

	if !rgx.MatchString(auc.email) {
		return errors.New("invalid email address")
	}

	return nil
}
