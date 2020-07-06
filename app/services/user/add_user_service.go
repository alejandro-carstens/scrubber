package user

import (
	"errors"
	"fmt"
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/user/contexts"
)

func NewAddUserService() *AddUserService {
	return &AddUserService{
		userRepository: repositories.NewUserRepository(),
	}
}

type AddUserService struct {
	userRepository *repositories.UserRepository
}

func (aus *AddUserService) Handle(context *contexts.AddUserContext) (*models.User, error) {
	users := []*models.User{}

	if err := aus.userRepository.FindWhere(map[string]interface{}{
		"email = ?": context.Email(),
	}, &users); err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return nil, errors.New("a user with the email provided already exists")
	}

	user := &models.User{
		Email:         context.Email(),
		EmailVerified: false,
		Name:          context.Name(),
		LastName:      context.LastName(),
	}

	user.FullName = fmt.Sprintf("%v %v", user.Name, user.LastName)

	return user, aus.userRepository.Create(user)
}
