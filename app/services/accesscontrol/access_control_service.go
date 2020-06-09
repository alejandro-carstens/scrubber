package accesscontrol

import (
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/jinzhu/gorm"
)

func NewAccessControlService() *AccessControlService {
	return &AccessControlService{
		accessControlRepository: repositories.NewAcessControlRepository(),
		userRepository:          repositories.NewUserRepository(),
	}
}

type AccessControlService struct {
	accessControlRepository *repositories.AccessControlRepository
	userRepository          *repositories.UserRepository
}

func (acs *AccessControlService) Handle(context *contexts.AccessControlContext) (*models.User, error) {
	accessControls := []*models.AccessControl{}

	if err := acs.accessControlRepository.FindWhere(map[string]interface{}{
		"user_id = ?": context.UserID(),
	}, &accessControls); err != nil {
		return nil, err
	}

	accessControlMap := context.AccessControlsMap()

	readIds := []uint64{}
	writeIds := []uint64{}
	noAccessIds := []uint64{}

	for _, accesssControl := range accessControls {
		entity, valid := accessControlMap[accesssControl.Action]

		if !valid {
			continue
		}

		delete(accessControlMap, accesssControl.Action)

		switch entity.Scope {
		case repositories.ACCESS_CONTROL_READ_SCOPE:
			readIds = append(readIds, accesssControl.ID)
			break
		case repositories.ACCESS_CONTROL_WRITE_SCOPE:
			writeIds = append(writeIds, accesssControl.ID)
			break
		case repositories.ACCESS_CONTROL_NO_ACCESS_SCOPE:
			noAccessIds = append(noAccessIds, accesssControl.ID)
			break
		}
	}

	inserts := []interface{}{}

	for _, accessControl := range accessControlMap {
		inserts = append(inserts, &models.AccessControl{
			UserID: context.UserID(),
			Action: accessControl.Action,
			Scope:  accessControl.Scope,
		})
	}

	if err := acs.accessControlRepository.DB().Transaction(func(tx *gorm.DB) error {
		accessControlRepository := acs.accessControlRepository.FromTx(tx)

		if len(readIds) > 0 {
			if _, err := accessControlRepository.UpdateWhere(map[string]interface{}{
				"id IN (?)": readIds,
			}, map[string]interface{}{
				"scope": repositories.ACCESS_CONTROL_READ_SCOPE,
			}); err != nil {
				return err
			}
		}

		if len(writeIds) > 0 {
			if _, err := accessControlRepository.UpdateWhere(map[string]interface{}{
				"id IN (?)": writeIds,
			}, map[string]interface{}{
				"scope": repositories.ACCESS_CONTROL_WRITE_SCOPE,
			}); err != nil {
				return err
			}
		}

		if len(noAccessIds) > 0 {
			if _, err := accessControlRepository.UpdateWhere(map[string]interface{}{
				"id IN (?)": noAccessIds,
			}, map[string]interface{}{
				"scope": repositories.ACCESS_CONTROL_NO_ACCESS_SCOPE,
			}); err != nil {
				return err
			}
		}

		if len(inserts) > 0 {
			return accessControlRepository.Insert(inserts...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	user := &models.User{}

	return user, acs.userRepository.Preload("AccessControls").Find(context.UserID(), user)
}
