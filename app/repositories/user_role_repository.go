package repositories

import (
	"scrubber/app/models"

	"github.com/jinzhu/gorm"
)

func NewUserRoleRepository() *UserRoleRepository {
	return repo(&models.UserRole{}, nil).(*UserRoleRepository)
}

type UserRoleRepository struct {
	repository
}

func (urr *UserRoleRepository) FromTx(tx *gorm.DB) *UserRoleRepository {
	return repo(&models.UserRole{}, tx).(*UserRoleRepository)
}
