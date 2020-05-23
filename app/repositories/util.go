package repositories

import (
	"github.com/alejandro-carstens/scrubber/app/models"
	rp "github.com/alejandro-carstens/scrubber/resourcepool"
)

func repo(model models.Modelable) Repositoryable {
	var repo Repositoryable

	switch model.Table() {
	case "users":
		repo = new(UserRepository)
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
