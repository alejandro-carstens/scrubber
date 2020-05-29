package contexts

import "errors"

func NewGoogleAuthContext(params map[string]interface{}) (*GoogleAuthContext, error) {
	context := &GoogleAuthContext{
		code: params["code"].(string),
	}

	return context.validate()
}

type GoogleAuthContext struct {
	code string
}

func (gac *GoogleAuthContext) Code() string {
	return gac.code
}

func (gac *GoogleAuthContext) validate() (*GoogleAuthContext, error) {
	if len(gac.code) == 0 {
		return nil, errors.New("code needs to be set")
	}

	return gac, nil
}
