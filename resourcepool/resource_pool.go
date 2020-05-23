package resourcepool

import (
	"context"
	"sync"

	"github.com/jinzhu/gorm"
)

var (
	rPool *ResourcePool
	once  sync.Once
)

func register() map[string]bootable {
	return map[string]bootable{
		"mysql": &mysql{},
	}
}

// Boot bootstraps the resources used by the application
func Boot(exclude ...string) error {
	var err error

	if rPool != nil {
		return err
	}

	once.Do(func() {
		rPool = &ResourcePool{
			context: context.Background(),
		}

		registry := register()

		for _, key := range exclude {
			delete(registry, key)
		}

		for _, resource := range registry {
			err = resource.boot(rPool)

			if err != nil {
				return
			}
		}
	})

	return err
}

// RPool returns an instance of *ResourcePool
func RPool() *ResourcePool {
	return rPool
}

// Context returns a context.Context
func Context() context.Context {
	return rPool.Context()
}

// MySQL returns a GORM MySQL database
func MySQL() *gorm.DB {
	return rPool.MySQL()
}

// ResourcePool represents the resources
// used by the application
type ResourcePool struct {
	context context.Context
	mysql   *gorm.DB
}

// Context returns a context.Context
func (rp *ResourcePool) Context() context.Context {
	return rp.context
}

// MySQL returns a GORM MySQL database
func (rp *ResourcePool) MySQL() *gorm.DB {
	return rp.mysql
}
