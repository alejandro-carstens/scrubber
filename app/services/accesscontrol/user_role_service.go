package accesscontrol

import (
	"errors"
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/jinzhu/gorm"
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

	roles, err := urs.getRoles(context.RoleIDs())

	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		if err := urs.deleteAllRolesForUser(user.ID); err != nil {
			return nil, err
		}

		return user, urs.userRepository.Preload("Roles.Permissions").Find(user.ID, user)
	}

	userRoles := []*models.UserRole{}

	if err := urs.userRoleRepository.FindWhere(map[string]interface{}{
		"user_id = ?": user.ID,
	}, &userRoles); err != nil {
		return nil, err
	}

	deletes := []uint64{}
	roleIdsMap := context.RoleIDsMap()

	for _, userRole := range userRoles {
		if _, valid := roleIdsMap[userRole.RoleID]; valid {
			delete(roleIdsMap, userRole.RoleID)

			continue
		}

		deletes = append(deletes, userRole.RoleID)
	}

	if err := urs.userRoleRepository.DB().Transaction(func(tx *gorm.DB) error {
		if len(deletes) > 0 {
			if _, err := urs.userRoleRepository.DeleteWhere(map[string]interface{}{
				"user_id = ?":    user.ID,
				"role_id IN (?)": deletes,
			}, &models.UserRole{}, true); err != nil {
				return err
			}
		}

		inserts := []interface{}{}

		for _, roleId := range roleIdsMap {
			inserts = append(inserts, &models.UserRole{UserID: context.UserID(), RoleID: roleId})
		}

		if len(inserts) > 0 {
			return urs.userRoleRepository.Insert(inserts...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return user, urs.userRepository.
		Preload("Roles.Permissions").Find(user.ID, user)
}

func (urs *UserRoleService) deleteAllRolesForUser(userId uint64) error {
	_, err := urs.userRoleRepository.DeleteWhere(map[string]interface{}{
		"user_id = ?": userId,
	}, &models.UserRole{}, true)

	return err
}

func (urs *UserRoleService) getRoles(roleIds []uint64) ([]*models.Role, error) {
	roles := []*models.Role{}

	if len(roleIds) == 0 {
		return roles, nil
	}

	if err := urs.roleRepository.FindWhere(map[string]interface{}{
		"id IN(?)": roleIds,
	}, &roles); err != nil {
		return nil, err
	}

	if len(roleIds) != len(roles) {
		return nil, errors.New("could not retrieve the specified roles")
	}

	return roles, nil
}
