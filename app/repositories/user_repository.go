package repositories

import "scrubber/app/models"

func UserRepo() *UserRepository {
	return repo(&models.User{}).(*UserRepository)
}

type UserRepository struct {
	repository
}
