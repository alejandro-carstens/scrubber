package user

import (
	"scrubber/app/models"
	"scrubber/app/repositories"
	"scrubber/app/services/user/contexts"

	"github.com/jinzhu/gorm"
)

func NewDeleteUserService() *DeleteUserService {
	return &DeleteUserService{
		userRepository:     repositories.NewUserRepository(),
		userRoleRepository: repositories.NewUserRoleRepository(),
	}
}

type DeleteUserService struct {
	userRepository     *repositories.UserRepository
	userRoleRepository *repositories.UserRoleRepository
}

func (dus *DeleteUserService) Handle(context *contexts.DeleteUserContext) error {
	user := &models.User{}

	if err := dus.userRepository.Find(context.UserID(), user); err != nil {
		return err
	}

	return dus.userRepository.DB().Transaction(func(tx *gorm.DB) error {
		if _, err := dus.userRoleRepository.FromTx(tx).DeleteWhere(map[string]interface{}{
			"user_id = ?": user.ID,
		}, &models.UserRole{}, false); err != nil {
			return err
		}

		_, err := dus.userRepository.FromTx(tx).Delete(user, false)

		return err
	})
}
