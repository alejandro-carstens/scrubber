package resourcepool

import (
	"fmt"
	"scrubber/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type mysql struct{}

func (m *mysql) boot(resourcePool *ResourcePool) error {
	configuration, err := config.GetFromEnvVars("mysql")

	if err != nil {
		return err
	}

	cnf := configuration.(*config.MySQL)

	mysql, err := gorm.Open(
		"mysql",
		fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
			cnf.Username,
			cnf.Password,
			cnf.Host,
			cnf.Port,
			cnf.Database,
		),
	)

	if err != nil {
		return err
	}

	if cnf.ConnMaxLifetime > 0 {
		mysql.DB().SetMaxOpenConns(cnf.ConnMaxLifetime)
	}

	if cnf.MaxIdleConnections > 0 {
		mysql.DB().SetMaxIdleConns(cnf.MaxIdleConnections)
	}

	if cnf.MaxOpenConnections > 0 {
		mysql.DB().SetMaxOpenConns(cnf.MaxOpenConnections)
	}

	resourcePool.mysql = mysql

	return nil
}
