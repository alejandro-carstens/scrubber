package repositories

import "scrubber/app/models"

func AccessControlRepo() *AccessControlRepository {
	return repo(&models.AccessControl{}).(*AccessControlRepository)
}

type AccessControlRepository struct {
	repository
}
