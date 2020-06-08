package accesscontrol

import (
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/accesscontrol/contexts"

	"github.com/jinzhu/gorm"
)

func NewAccessControlService() *AccessControlService {
	return &AccessControlService{
		repository: repositories.NewAcessControlRepository(),
	}
}

type AccessControlService struct {
	repository *repositories.AccessControlRepository
}

func (acs *AccessControlService) Handle(context *contexts.AccessControlContext) ([]*models.AccessControl, error) {
	accessControls := []*models.AccessControl{}

	if err := acs.repository.FindWhere(map[string]interface{}{
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

	if err := acs.repository.DB().Transaction(func(tx *gorm.DB) error {
		repository := acs.repository.FromTx(tx)

		if len(readIds) > 0 {
			if _, err := repository.UpdateWhere(map[string]interface{}{
				"id IN (?)": readIds,
			}, map[string]interface{}{
				"scope": repositories.ACCESS_CONTROL_READ_SCOPE,
			}); err != nil {
				return err
			}
		}

		if len(writeIds) > 0 {
			if _, err := repository.UpdateWhere(map[string]interface{}{
				"id IN (?)": writeIds,
			}, map[string]interface{}{
				"scope": repositories.ACCESS_CONTROL_WRITE_SCOPE,
			}); err != nil {
				return err
			}
		}

		if len(noAccessIds) > 0 {
			if _, err := repository.UpdateWhere(map[string]interface{}{
				"id IN (?)": noAccessIds,
			}, map[string]interface{}{
				"scope": repositories.ACCESS_CONTROL_NO_ACCESS_SCOPE,
			}); err != nil {
				return err
			}
		}

		if len(inserts) > 0 {
			return repository.Insert(inserts...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	accessControls = []*models.AccessControl{}

	err := acs.repository.FindWhere(map[string]interface{}{
		"user_id = ?": context.UserID(),
	}, &accessControls)

	return accessControls, err
}
