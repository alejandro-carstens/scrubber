package repositories

import (
	"scrubber/app/models"

	"github.com/jinzhu/gorm"
)

type Repositoryable interface {
	DB() *gorm.DB

	Tx(tx *gorm.DB) Repositoryable

	Unscoped() Repositoryable

	Preload(relation string, conditions ...interface{}) Repositoryable

	Find(id uint64, accountId uint64, dest interface{}) error

	FindWhere(params map[string]interface{}, dest interface{}) error

	Create(model models.Modelable) error

	Insert(inserts ...interface{}) error

	Update(accountId uint64, model models.Modelable) error

	UpdateWhere(params map[string]interface{}, updates map[string]interface{}) (int64, error)

	DeleteWhere(params map[string]interface{}, model models.Modelable, hard bool) (int64, error)

	Init(model models.Modelable, db *gorm.DB) Repositoryable

	ReleaseTx()
}
