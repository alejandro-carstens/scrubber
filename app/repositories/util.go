package repositories

import (
	"scrubber/app/models"
	rp "scrubber/resourcepool"

	"github.com/jinzhu/gorm"
)

func repo(model models.Modelable, db *gorm.DB) Repositoryable {
	var repository Repositoryable

	switch model.Table() {
	case "users":
		repository = new(UserRepository)
	case "access_controls":
		repository = new(AccessControlRepository)
	}

	if db == nil {
		db = rp.MySQL()
	}

	repository.Init(model, db)

	return repository
}

func copyMap(m map[string][]interface{}) map[string][]interface{} {
	x := map[string][]interface{}{}

	for k, v := range m {
		x[k] = v
	}

	return x
}

func inStringSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}
