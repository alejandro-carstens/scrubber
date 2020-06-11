package repositories

import (
	"scrubber/app/models"

	"github.com/jinzhu/gorm"
)

func NewPermissionRepository() *PermissionRepository {
	return repo(&models.Permission{}, nil).(*PermissionRepository)
}

type PermissionRepository struct {
	repository
}

func (ur *PermissionRepository) FromTx(tx *gorm.DB) *PermissionRepository {
	return repo(&models.Permission{}, tx).(*PermissionRepository)
}
