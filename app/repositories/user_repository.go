package repositories

import (
	"scrubber/app/models"

	"github.com/jinzhu/gorm"
)

func NewUserRepository() *UserRepository {
	return repo(&models.User{}, nil).(*UserRepository)
}

type UserRepository struct {
	repository
}

func (ur *UserRepository) FromTx(tx *gorm.DB) *UserRepository {
	return repo(&models.User{}, tx).(*UserRepository)
}
