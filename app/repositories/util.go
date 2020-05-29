package repositories

import (
	"scrubber/app/models"
	rp "scrubber/resourcepool"
)

func repo(model models.Modelable) Repositoryable {
	var repo Repositoryable

	switch model.Table() {
	case "users":
		repo = new(UserRepository)
	case "access_controls":
		repo = new(AccessControlRepository)
	}

	repo.Init(model, rp.MySQL())

	return repo
}

func copyMap(m map[string][]interface{}) map[string][]interface{} {
	x := map[string][]interface{}{}

	for k, v := range m {
		x[k] = v
	}

	return x
}
