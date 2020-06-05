package repositories

import (
	"scrubber/app/models"

	"github.com/jinzhu/gorm"
)

const ACCESS_CONTROL_ALL_ACTIONS string = "all"

const ACCESS_CONTROL_READ_SCOPE string = "read"
const ACCESS_CONTROL_WRITE_SCOPE string = "write"

func NewAcessControlRepository() *AccessControlRepository {
	return repo(&models.AccessControl{}, nil).(*AccessControlRepository)
}

type AccessControlRepository struct {
	repository
}

func (acr *AccessControlRepository) FromTx(tx *gorm.DB) *AccessControlRepository {
	return repo(&models.AccessControl{}, tx).(*AccessControlRepository)
}
