package resourcepool

import (
	"context"
	"errors"
	"sync"

	logs "scrubber/logger"

	"github.com/alejandro-carstens/golastic"
	"github.com/jinzhu/gorm"
)

var (
	rPool *ResourcePool
	once  sync.Once
)

func register() map[string]bootable {
	return map[string]bootable{
		"mysql":         &mysql{},
		"elasticsearch": &elasticsearch{},
		"logger":        &logger{},
	}
}

func IsBooted() bool {
	return rPool != nil
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

func BootResource(name string) error {
	if !IsBooted() {
		return errors.New("a resource pool needs to be booted")
	}

	resource, valid := register()[name]

	if !valid {
		return errors.New("invalid resource specified")
	}

	return resource.boot(rPool)
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

// Elasticsearch returns an Elasticsearch connection
func Elasticsearch() *golastic.Connection {
	return rPool.Elasticsearch()
}

func Logger() *logs.Logger {
	return rPool.Logger()
}

// ResourcePool represents the resources
// used by the application
type ResourcePool struct {
	context       context.Context
	mysql         *gorm.DB
	elasticsearch *golastic.Connection
	logger        *logs.Logger
}

// Context returns a context.Context
func (rp *ResourcePool) Context() context.Context {
	return rp.context
}

// MySQL returns a GORM MySQL database
func (rp *ResourcePool) MySQL() *gorm.DB {
	return rp.mysql
}

// Elasticsearch returns an Elasticsearch connection
func (rp *ResourcePool) Elasticsearch() *golastic.Connection {
	return rp.elasticsearch
}

func (rp *ResourcePool) Logger() *logs.Logger {
	return rp.logger
}
