package repositories

import "github.com/alejandro-carstens/scrubber/app/models"

func UserRepo() *UserRepository {
	return repo(&models.User{}).(*UserRepository)
}

type UserRepository struct {
	repository
}
