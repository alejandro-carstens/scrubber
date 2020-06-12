package accesscontrol

import (
	"errors"
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/accesscontrol/contexts"
)

func NewUserRoleService() *UserRoleService {
	return &UserRoleService{
		userRepository:     repositories.NewUserRepository(),
		roleRepository:     repositories.NewRoleRepository(),
		userRoleRepository: repositories.NewUserRoleRepository(),
	}
}

type UserRoleService struct {
	userRepository     *repositories.UserRepository
	roleRepository     *repositories.RoleRepository
	userRoleRepository *repositories.UserRoleRepository
}

func (urs *UserRoleService) Handle(context *contexts.UserRoleContext) (*models.User, error) {
	user := &models.User{}

	if err := urs.userRepository.Find(context.UserID(), user); err != nil {
		return nil, err
	}

	role := &models.Role{}

	if err := urs.roleRepository.Find(context.RoleID(), role); err != nil {
		return nil, err
	}

	userRoles := []*models.UserRole{}

	if err := urs.userRoleRepository.FindWhere(map[string]interface{}{
		"user_id = ?": user.ID,
	}, &userRoles); err != nil {
		return nil, err
	}

	if len(userRoles) > 1 {
		return nil, errors.New("There can be at most one role per user")
	}

	if len(userRoles) == 0 {
		if err := urs.userRoleRepository.Create(&models.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}); err != nil {
			return nil, err
		}
	} else {
		userRole := userRoles[0]
		userRole.RoleID = role.ID

		if err := urs.userRoleRepository.Update(userRole); err != nil {
			return nil, err
		}
	}

	return user, urs.userRepository.Preload("Roles.Permissions").Find(user.ID, user)
}
