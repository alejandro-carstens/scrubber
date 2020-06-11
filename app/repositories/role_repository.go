package repositories

import (
	"scrubber/app/models"

	"github.com/jinzhu/gorm"
)

func NewRoleRepository() *RoleRepository {
	return repo(&models.Role{}, nil).(*RoleRepository)
}

type RoleRepository struct {
	repository
}

func (rr *RoleRepository) FromTx(tx *gorm.DB) *RoleRepository {
	return repo(&models.Role{}, tx).(*RoleRepository)
}
