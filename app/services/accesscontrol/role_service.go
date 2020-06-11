package accesscontrol

import (
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/jinzhu/gorm"
)

func NewRoleService() *RoleService {
	return &RoleService{
		permissionRepository: repositories.NewPermissionRepository(),
		roleRepository:       repositories.NewRoleRepository(),
	}
}

type RoleService struct {
	permissionRepository *repositories.PermissionRepository
	roleRepository       *repositories.RoleRepository
}

func (rs *RoleService) Handle(context *contexts.RoleContext) (*models.Role, error) {
	permissions := []*models.Permission{}

	if context.RoleID() > 0 {
		if err := rs.permissionRepository.FindWhere(map[string]interface{}{
			"role_id = ?": context.RoleID(),
		}, &permissions); err != nil {
			return nil, err
		}
	}

	readIds := []uint64{}
	writeIds := []uint64{}
	noAccessIds := []uint64{}
	permissionMap := context.PermissionMap()

	for _, permission := range permissions {
		entity, valid := permissionMap[permission.Action]

		if !valid {
			continue
		}

		if entity.Scope == repositories.READ_SCOPE {
			readIds = append(readIds, permission.ID)
		}

		if entity.Scope == repositories.WRITE_SCOPE {
			writeIds = append(writeIds, permission.ID)
		}

		if entity.Scope == repositories.NO_ACCESS_SCOPE {
			noAccessIds = append(noAccessIds, permission.ID)
		}

		delete(permissionMap, permission.Action)
	}

	role := &models.Role{Name: context.Name()}

	if err := rs.permissionRepository.DB().Transaction(func(tx *gorm.DB) error {
		permissionRepository := rs.permissionRepository.FromTx(tx)

		if context.RoleID() == 0 {
			if err := rs.roleRepository.FromTx(tx).Create(role); err != nil {
				return err
			}
		} else {
			role.ID = context.RoleID()

			if err := rs.roleRepository.Update(role); err != nil {
				return err
			}
		}

		if len(readIds) > 0 {
			if _, err := permissionRepository.UpdateWhere(
				rs.params(readIds, repositories.READ_SCOPE),
			); err != nil {
				return err
			}
		}

		if len(writeIds) > 0 {
			if _, err := permissionRepository.UpdateWhere(
				rs.params(writeIds, repositories.WRITE_SCOPE),
			); err != nil {
				return err
			}
		}

		if len(noAccessIds) > 0 {
			if _, err := permissionRepository.UpdateWhere(
				rs.params(noAccessIds, repositories.NO_ACCESS_SCOPE),
			); err != nil {
				return err
			}
		}

		if len(permissionMap) > 0 {
			return permissionRepository.Insert(rs.prepareInserts(role.ID, permissionMap)...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return role, rs.roleRepository.Preload("Permissions").Find(role.ID, role)
}

func (rs *RoleService) prepareInserts(roleId uint64, permissionMap map[string]*contexts.PermissionEntity) []interface{} {
	inserts := []interface{}{}

	for _, permission := range permissionMap {
		inserts = append(inserts, &models.Permission{
			RoleID: roleId,
			Action: permission.Action,
			Scope:  permission.Scope,
		})
	}

	return inserts
}

func (rs *RoleService) params(ids []uint64, scope string) (map[string]interface{}, map[string]interface{}) {
	return map[string]interface{}{"id IN (?)": ids}, map[string]interface{}{"scope": scope}
}
