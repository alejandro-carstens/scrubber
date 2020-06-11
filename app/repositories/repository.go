package repositories

import (
	"scrubber/app/models"
	"sync"

	"github.com/jinzhu/gorm"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"
)

const CHUNK_SIZE int = 2500
const LIMIT int = 10000

type repository struct {
	db       *gorm.DB
	tx       *gorm.DB
	model    models.Modelable
	preloads map[string][]interface{}
	unscoped bool
}

func (r *repository) DB() *gorm.DB {
	return r.db
}

func (r *repository) Init(model models.Modelable, db *gorm.DB) Repositoryable {
	r.db = db
	r.model = model

	return r
}

func (r *repository) Unscoped() Repositoryable {
	r = r.clone()

	r.unscoped = true

	return r
}

func (r *repository) Preload(relation string, conditions ...interface{}) Repositoryable {
	r = r.clone()

	if r.preloads == nil {
		r.preloads = map[string][]interface{}{}
	}

	r.preloads[relation] = conditions

	return r
}

func (r *repository) FindWhere(params map[string]interface{}, dest interface{}) error {
	query := r.connection().Table(r.model.Table()).LogMode(true)

	if r.unscoped {
		query = query.Unscoped()
	}

	if len(r.preloads) > 0 {
		for relation, conditions := range r.preloads {
			query = query.Preload(relation, conditions...)
		}
	}

	for condition, value := range params {
		query = query.Where(condition, value)
	}

	query = query.Limit(LIMIT).Find(dest)

	return query.Error
}

func (r *repository) Find(id uint64, dest interface{}) error {
	query := r.connection().Table(r.model.Table()).LogMode(true)

	if r.unscoped {
		query = query.Unscoped()
	}

	if len(r.preloads) > 0 {
		for relation, conditions := range r.preloads {
			query = query.Preload(relation, conditions...)
		}
	}

	return query.First(dest, id).Error
}

func (r *repository) Create(model models.Modelable) error {
	res := r.connection().Table(r.model.Table()).LogMode(true).Create(model)

	return res.Error
}

func (r *repository) Insert(inserts ...interface{}) error {
	return gormbulk.BulkInsert(
		r.connection().Table(r.model.Table()).LogMode(true),
		inserts,
		CHUNK_SIZE,
	)
}

func (r *repository) UpdateWhere(params map[string]interface{}, updates map[string]interface{}) (int64, error) {
	query := r.connection().Table(r.model.Table()).LogMode(true)

	if r.unscoped {
		query.Unscoped()
	}

	for condition, value := range params {
		query = query.Where(condition, value)
	}

	res := query.Limit(LIMIT).Updates(updates)

	return res.RowsAffected, res.Error
}

func (r *repository) Update(model models.Modelable) error {
	res := r.connection().Table(r.model.Table()).LogMode(true).Save(model)

	return res.Error
}

func (r *repository) DeleteWhere(params map[string]interface{}, model models.Modelable, hard bool) (int64, error) {
	query := r.connection().Table(r.model.Table()).LogMode(true)

	if r.unscoped {
		query.Unscoped()
	}

	for condition, value := range params {
		query = query.Where(condition, value)
	}

	res := query.Limit(LIMIT).Delete(model)

	return res.RowsAffected, res.Error
}

func (r *repository) QueryByContext(context *QueryContext, dest interface{}) (*queryMeta, error) {
	query := r.buildQueryFromContext(context)
	metaQuery := query

	var limit int

	if context.Limit > 0 {
		limit = context.Limit
	} else {
		limit = LIMIT
	}

	query = query.Limit(limit)
	query = query.Offset(context.Offset)

	var wg sync.WaitGroup

	wg.Add(2)

	var total int

	go func() {
		defer wg.Done()

		metaQuery = metaQuery.Count(&total)
	}()

	go func() {
		defer wg.Done()

		query = query.Find(dest)
	}()

	wg.Wait()

	if metaQuery.Error != nil {
		return nil, metaQuery.Error
	}

	if query.Error != nil {
		return nil, query.Error
	}

	return buildQueryMeta(limit, context.Offset, total), nil
}

func (r *repository) buildQueryFromContext(context *QueryContext) *gorm.DB {
	query := r.connection().Table(r.model.Table()).LogMode(true)

	if r.unscoped {
		query = query.Unscoped()
	}

	for _, include := range context.Includes {
		query = query.Preload(include)
	}

	for _, where := range context.Wheres {
		query = query.Where(where.Prepare())
	}

	for _, whereNull := range context.WhereNulls {
		query = query.Where(whereNull.Prepare())
	}

	for _, whereNotNull := range context.WhereNotNulls {
		query = query.Where(whereNotNull.Prepare())
	}

	for _, orWhere := range context.OrWheres {
		query = query.Or(orWhere.Prepare())
	}

	for _, whereIn := range context.WhereIns {
		query = query.Where(whereIn.Prepare())
	}

	for _, whereNotIn := range context.WhereNotIns {
		query = query.Not(whereNotIn.Prepare())
	}

	return query
}

func (r *repository) clone() *repository {
	clone := *r

	if len(r.preloads) > 0 {
		clone.preloads = copyMap(r.preloads)
	}

	if r.unscoped {
		clone.unscoped = true
	}

	return &clone
}

func (r *repository) connection() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}

	return r.db
}
