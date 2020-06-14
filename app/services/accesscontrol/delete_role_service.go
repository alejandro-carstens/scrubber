package accesscontrol

import (
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/jinzhu/gorm"
)

func NewDeleteRoleService() *DeleteRoleService {
	return &DeleteRoleService{
		roleRepository:       repositories.NewRoleRepository(),
		userRoleRespository:  repositories.NewUserRoleRepository(),
		permissionRepository: repositories.NewPermissionRepository(),
	}
}

type DeleteRoleService struct {
	roleRepository       *repositories.RoleRepository
	userRoleRespository  *repositories.UserRoleRepository
	permissionRepository *repositories.PermissionRepository
}

func (drs *DeleteRoleService) Handle(context *contexts.DeleteRoleContext) error {
	role := &models.Role{}

	if err := drs.roleRepository.Find(context.RoleID(), role); err != nil {
		return err
	}

	return drs.roleRepository.DB().Transaction(func(tx *gorm.DB) error {
		if _, err := drs.permissionRepository.FromTx(tx).DeleteWhere(map[string]interface{}{
			"role_id = ?": role.ID,
		}, &models.Permission{}, false); err != nil {
			return err
		}

		if _, err := drs.userRoleRespository.FromTx(tx).DeleteWhere(map[string]interface{}{
			"role_id = ?": role.ID,
		}, &models.UserRole{}, false); err != nil {
			return err
		}

		_, err := drs.roleRepository.FromTx(tx).Delete(role, false)

		return err
	})
}
